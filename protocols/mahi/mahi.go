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

type Mahi struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewMahi(logger *util.Logger) *Mahi {
	return &Mahi{
		logger: logger,
	}
}

func (ba *Mahi) CopyConsensus(nodes []*common.Node) error {
	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("num_replicas not found in options")
	}

	num_replicas_int, err := strconv.ParseInt(num_replicas, 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}
	if num_replicas_int > int64(len(nodes)) {
		panic("Not enough nodes to deploy mahi")
	}

	nodes[0].ExecCmd(fmt.Sprintf("sudo rm -r async-mystecity; git clone https://github.com/PasinduTennage/async-mystecity; cd async-mystecity; git checkout consensus-rework; sudo apt-get install -y libfontconfig1-dev; source %v.cargo/env; cargo build", nodes[0].HomeDir))

	println("Cloned the mahi repository and built the binary")

	nodes[0].Get_Load(fmt.Sprintf("%vasync-mystecity/target/debug/mysticeti", nodes[0].HomeDir), "protocols/mahi/assets/")

	println("Copied the mahi binary to the controller machine")

	var wg sync.WaitGroup
	wg.Add(int(num_replicas_int))

	for j := int64(0); j < num_replicas_int; j++ {
		go func(i int) {
			nodes[i].Put_Load("protocols/mahi/assets/mysticeti", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/mahi/assets/config-rewrite.py", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			wg.Done()
		}(int(j))
	}
	wg.Wait()
	fmt.Print("Copied the mahi binary to all the nodes\n")
	return nil
}

func (ba *Mahi) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
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

	fmt.Print("Killed all the replicas and clients\n")

	load, ok := ba.options.Option["load"]
	if !ok {
		panic("load not found in options")
	}

	wave_length, ok := ba.options.Option["wave_length"]
	if !ok {
		panic("wave_length not found in options")
	}

	number_of_leaders, ok := ba.options.Option["number_of_leaders"]
	if !ok {
		panic("number_of_leaders not found in options")
	}

	enable_pipelining, ok := ba.options.Option["enable_pipelining"]
	if !ok {
		panic("enable_pipelining not found in options")
	}

	enable_synchronizer, ok := ba.options.Option["enable_synchronizer"]
	if !ok {
		panic("enable_synchronizer not found in options")
	}

	transaction_size, ok := ba.options.Option["transaction_size"]
	if !ok {
		panic("transaction_size not found in options")
	}

	sshCmd := exec.Command("python3", []string{"protocols/mahi/assets/genrate-configs.py", "--wave_length", wave_length, "--number_of_leaders", number_of_leaders, "--enable_pipelining", enable_pipelining, "--consensus_only", "true", "--enable_synchronizer", enable_synchronizer, "--initial_delay_secs", "5", "--initial_delay_nanos", "0", "--load", load, "--transaction_size", transaction_size, "--output_dir", "protocols/mahi/assets/"}...)
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
			nodes[j].Put_Load("protocols/mahi/assets/client-parameters.yml", fmt.Sprintf("%vbench/", nodes[j].HomeDir))
			nodes[j].Put_Load("protocols/mahi/assets/node-parameters.yml", fmt.Sprintf("%vbench/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("rm %vbench/storage-%v/wal", nodes[j].HomeDir, j))
			nodes[j].ExecCmd(fmt.Sprintf("./bench/mysticeti benchmark-genesis --ips %v --working-directory %v --node-parameters-path %vnode-parameters.yml", ip_string, nodes[j].HomeDir+"bench/", nodes[j].HomeDir+"bench/"))
			nodes[j].ExecCmd(fmt.Sprintf("python3 %vbench/config-rewrite.py %vbench/public-config.yaml", nodes[j].HomeDir, nodes[j].HomeDir))
			wg1.Done()
		}(i)
	}
	wg1.Wait()

	println("Generated the node private keys")

	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].ExecCmd(fmt.Sprintf("./bench/mysticeti run --authority %v --committee-path %vbench/committee.yaml --public-config-path %vbench/public-config.yaml --private-config-path %vbench/private-config-%v.yaml --client-parameters-path %vbench/client-parameters.yml", j, nodes[j].HomeDir, nodes[j].HomeDir, nodes[j].HomeDir, j, nodes[j].HomeDir))
		}(i)
	}
	time.Sleep(5 * time.Second)

	println("Started the replicas")

	bootstrap_complete <- true

	time.Sleep(time.Duration(duration) * time.Second)

	var wg2 sync.WaitGroup
	wg2.Add(int(num_replicas))
	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].ExecCmd("pkill -KILL -f mysticeti")
			wg2.Done()
		}(i)
	}
	wg2.Wait()

	println("Killed all the replicas")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Error getting home directory:" + err.Error())
	}

	cmd := exec.Command("rm", []string{"-r", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		print("Error while deleting logs/ " + err.Error() + " " + string(output) + "\n")
	} else {
		print("deleted local logs/ successfully\n" + string(output) + "\n")
	}

	cmd = exec.Command("mkdir", []string{"-p", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		panic("Error while creating logs/ " + err.Error() + " " + string(output) + "\n")
	} else {
		print("created logs/ successfully\n" + string(output) + "\n")
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

	println("Downloaded the client logs")

	command := "protocols/mahi/assets/performance_graph.py"
	outputs := []string{}
	for j := 0; j < int(num_replicas); j++ {
		logFile := filepath.Join(homeDir, fmt.Sprintf("toture-testing-consensus/logs/client-times-%v.txt", j))

		sshCmd = exec.Command("python3", []string{command, "mahi-" + strconv.Itoa(j), logFile}...)
		output, err = sshCmd.CombinedOutput()
		if err != nil {
			print("Error while generating performance graph " + err.Error() + " " + string(output) + "\n")
		} else {
			print("Mahi Performance graph generated successfully\n" + string(output) + "\n")
			outputs = append(outputs, string(output))
		}
	}

	fmt.Printf("Mahi Mahi Performance:\n %v\n", outputs)
	result <- ba.getPerformance(outputs)

}

func (ba *Mahi) ExtractOptions(path string) protocols.ConsensusOptions {
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

	fmt.Printf("Mahi Mahi options:\n %v\n", options.Option)

	ba.options = options
	return options
}

func (ba *Mahi) getPerformance(outputs []string) util.Performance {
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

	return p
}
