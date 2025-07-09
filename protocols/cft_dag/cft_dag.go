package consensus

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"toture-test/consenbench/common"
	"toture-test/protocols"
	"toture-test/util"
)

type CFT_DAG struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewCFT_DAG(logger *util.Logger) *CFT_DAG {
	return &CFT_DAG{
		logger: logger,
	}
}

func (ba *CFT_DAG) CopyConsensus(nodes []*common.Node) error {
	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("num_replicas not found in options")
	}

	num_replicas_int, err := strconv.ParseInt(num_replicas, 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}
	if num_replicas_int > int64(len(nodes)) {
		panic("Not enough nodes to deploy CFT_DAG")
	}

	nodes[0].ExecCmd(fmt.Sprintf("sudo rm -r mysticeti; git clone https://github.com/asonnino/mysticeti; cd mysticeti; git checkout cft; sudo apt-get install -y libfontconfig1-dev; source %v.cargo/env; cargo build --release", nodes[0].HomeDir))

	ba.logger.Debug(fmt.Sprintf("Cloned the CFT_DAG repository and built the binary"), 0)

	nodes[0].Get_Load(fmt.Sprintf("%vmysticeti/target/release/mysticeti", nodes[0].HomeDir), "protocols/cft_dag/assets/")

	ba.logger.Debug(fmt.Sprintf("Copied the CFT_DAG binary to the controller machine"), 0)

	var wg sync.WaitGroup
	wg.Add(int(num_replicas_int))

	for j := int64(0); j < num_replicas_int; j++ {
		go func(i int) {
			nodes[i].Put_Load("protocols/cft_dag/assets/mysticeti", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/cft_dag/assets/config-rewrite.py", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			wg.Done()
		}(int(j))
	}
	wg.Wait()
	ba.logger.Debug(fmt.Sprintf("Copied the CFT_DAG binary to all the nodes\n"), 0)
	return nil
}

func (ba *CFT_DAG) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {

	data, err := os.ReadFile("params.param")
	if err != nil {
		panic(err.Error())
	}

	parts := strings.Fields(string(data))
	if len(parts) != 2 {
		panic("Expected exactly 2 values in params.param")
	}

	param_load := parts[0]

	param_size := parts[1]

	num_replicas, err := strconv.ParseInt(ba.options.Option["num_replicas"], 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}

	var wg sync.WaitGroup
	wg.Add(int(num_replicas))
	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].ExecCmd("pkill -KILL -f mysticeti")
			nodes[j].ExecCmd(fmt.Sprintf("rm -r %vbench/logs/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vbench/logs/", nodes[j].HomeDir))
			wg.Done()
		}(i)
	}
	wg.Wait()

	ba.logger.Debug(fmt.Sprintf("Killed all the replicas\n"), 0)

	int_load, _ := strconv.Atoi(param_load)
	load := strconv.Itoa(int_load / int(num_replicas))

	wave_length := "2"

	number_of_leaders, ok := ba.options.Option["number_of_leaders"]
	if !ok {
		panic("number_of_leaders not found in options")
	}

	enable_pipelining := "true"

	enable_synchronizer := "true"

	transaction_size := param_size

	sshCmd := exec.Command("python3", []string{"protocols/cft_dag/assets/genrate-configs.py", "--wave_length", wave_length, "--number_of_leaders", number_of_leaders, "--enable_pipelining", enable_pipelining, "--consensus_only", "true", "--enable_synchronizer", enable_synchronizer, "--initial_delay_secs", "5", "--initial_delay_nanos", "0", "--load", load, "--transaction_size", transaction_size, "--output_dir", "protocols/cft_dag/assets/"}...)
	output, err := sshCmd.CombinedOutput()
	if err != nil {
		panic("Error while running config-generate.py " + err.Error() + " " + string(output))
	}

	ip_string := ""
	for i := 0; i < int(num_replicas); i++ {
		ip_string = ip_string + nodes[i].Ip + " "
	}

	var wg1 sync.WaitGroup
	wg1.Add(int(num_replicas))
	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].ExecCmd("pkill -KILL -f mysticeti")
			nodes[j].ExecCmd(fmt.Sprintf("rm %vclient-times-%v.txt", nodes[j].HomeDir, j))
			nodes[j].Put_Load("protocols/cft_dag/assets/client-parameters.yml", fmt.Sprintf("%vbench/", nodes[j].HomeDir))
			nodes[j].Put_Load("protocols/cft_dag/assets/node-parameters.yml", fmt.Sprintf("%vbench/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("rm %vbench/storage-%v/wal", nodes[j].HomeDir, j))
			nodes[j].ExecCmd(fmt.Sprintf("./bench/mysticeti benchmark-genesis --ips %v --working-directory %v --node-parameters-path %vnode-parameters.yml", ip_string, nodes[j].HomeDir+"bench/", nodes[j].HomeDir+"bench/"))
			nodes[j].ExecCmd(fmt.Sprintf("python3 %vbench/config-rewrite.py %vbench/public-config.yaml", nodes[j].HomeDir, nodes[j].HomeDir))
			wg1.Done()
		}(i)
	}
	wg1.Wait()

	ba.logger.Debug(fmt.Sprintf("Generated the node private keys"), 0)

	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].ExecCmd(fmt.Sprintf("./bench/mysticeti run --authority %v --committee-path %vbench/committee.yaml --public-config-path %vbench/public-config.yaml --private-config-path %vbench/private-config-%v.yaml --client-parameters-path %vbench/client-parameters.yml", j, nodes[j].HomeDir, nodes[j].HomeDir, nodes[j].HomeDir, j, nodes[j].HomeDir))
		}(i)
	}
	time.Sleep(45 * time.Second)

	ba.logger.Debug(fmt.Sprintf("Started the replicas"), 0)

	bootstrap_complete <- true

	time.Sleep(time.Duration(duration*3) * time.Second)

	var wg2 sync.WaitGroup
	wg2.Add(int(num_replicas))
	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].ExecCmd("pkill -KILL -f mysticeti")
			wg2.Done()
		}(i)
	}
	wg2.Wait()

	ba.logger.Debug(fmt.Sprintf("Killed all the replicas"), 0)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Error getting home directory:" + err.Error())
	}

	cmd := exec.Command("rm", []string{"-r", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug(fmt.Sprintf("error while deleting logs/ "+err.Error()+" "+string(output)+"\n"), 0)
	} else {
		ba.logger.Debug(fmt.Sprintf("deleted local logs/ successfully\n"+string(output)+"\n"), 0)
	}

	cmd = exec.Command("mkdir", []string{"-p", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		panic("Error while creating logs/ " + err.Error() + " " + string(output) + "\n")
	} else {
		ba.logger.Debug(fmt.Sprintf("created logs/ successfully\n"+string(output)+"\n"), 0)
	}

	var wg3 sync.WaitGroup
	wg3.Add(int(num_replicas))
	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].Get_Load(fmt.Sprintf("%vclient-times-%v.txt", nodes[j].HomeDir, j), fmt.Sprintf("logs/"))
			wg3.Done()
		}(i)
	}
	wg3.Wait()

	ba.logger.Debug(fmt.Sprintf("Downloaded the client logs"), 0)

	command := "protocols/cft_dag/assets/performance_graph.py"
	outputs := []string{}
	logFiles := []string{}
	for j := 0; j < int(num_replicas); j++ {
		logFile := filepath.Join(homeDir, fmt.Sprintf("toture-testing-consensus/logs/client-times-%v.txt", j))
		logFiles = append(logFiles, logFile)

		sshCmd = exec.Command("python3", []string{command, "CFT_DAG-" + strconv.Itoa(j), strconv.Itoa(duration), logFile}...)
		output, err = sshCmd.CombinedOutput()
		if err != nil {
			ba.logger.Debug(fmt.Sprintf("Error while generating performance graph "+err.Error()+" "+string(output)+"\n"), 0)
		} else {
			ba.logger.Debug(fmt.Sprintf("CFT_DAG Performance graph generated successfully\n"+string(output)+"\n"), 0)
			outputs = append(outputs, string(output))
		}
	}

	sshCmd = exec.Command("python3", append([]string{command, "CFT_DAG", strconv.Itoa(duration)}, logFiles...)...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug(fmt.Sprintf("Error while generating performance graph "+err.Error()+" "+string(output)+"\n"), 0)
	} else {
		ba.logger.Debug(fmt.Sprintf("CFT_DAG Performance graph generated successfully\n"+string(output)+"\n"), 0)
	}

	ba.logger.Debug(fmt.Sprintf("CFT_DAG Performance:\n %v\n", outputs), 0)
	result <- ba.getPerformance(outputs)

}

func (ba *CFT_DAG) ExtractOptions(path string) protocols.ConsensusOptions {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic("Error unmarshalling YAML " + err.Error())
	}

	options := protocols.ConsensusOptions{Option: make(map[string]string)}
	for key, value := range config {
		options.Option[key] = fmt.Sprintf("%v", value)
	}

	ba.logger.Debug(fmt.Sprintf("CFT_DAG options:\n %v\n", options.Option), 0)

	ba.options = options
	return options
}

func (ba *CFT_DAG) getPerformance(outputs []string) util.Performance {
	p := util.Performance{
		Option: make(map[string]string),
	}
	sum_tx := 0
	sum_lat := 0
	for i := 0; i < len(outputs); i++ {
		tx, err := strconv.ParseFloat(strings.Split(outputs[i], " ")[0], 64)
		if err != nil {
			panic(err.Error() + " while parsing tx")
		}
		lat, err := strconv.ParseFloat(strings.Split(outputs[i], " ")[1], 64)
		if err != nil {
			panic(err.Error() + " while parsing lat")
		}
		sum_tx += int(tx)
		sum_lat += int(lat)
	}
	p.Option["throughput"] = fmt.Sprintf("%v requests per second", sum_tx/len(outputs))
	p.Option["average latency"] = fmt.Sprintf("%v ms", sum_lat/len(outputs))

	fmt.Printf("%v,%v,%v,", sum_tx/len(outputs), sum_lat/len(outputs), 0)

	return p
}
