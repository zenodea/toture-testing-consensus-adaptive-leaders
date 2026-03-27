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

type Hammerhead2 struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewHammerhead2(logger *util.Logger) *Hammerhead2 {
	return &Hammerhead2{
		logger: logger,
	}
}

func (ba *Hammerhead2) CopyConsensus(nodes []*common.Node) error {
	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("num_replicas not found in options")
	}

	num_replicas_int, err := strconv.ParseInt(num_replicas, 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}
	if num_replicas_int > int64(len(nodes)) {
		panic("Not enough nodes to deploy hammerhead2")
	}

	nodes[0].ExecCmd(fmt.Sprintf("sudo rm -r hammerhead_2.0; git clone -b feat/leader_scheduling_v2 https://github.com/zenodea/hammerhead_2.0; cd hammerhead_2.0 ; sudo apt-get install -y libfontconfig1-dev; source %v.cargo/env; cargo build --release", nodes[0].HomeDir))

	ba.logger.Debug(fmt.Sprintf("Cloned the hammerhead2 repository and built the binary"), 0)

	nodes[0].Get_Load(fmt.Sprintf("%vhammerhead_2.0/target/release/mysticeti", nodes[0].HomeDir), "protocols/hammerhead2/assets/")

	ba.logger.Debug(fmt.Sprintf("Copied the hammerhead2 binary to the controller machine"), 0)

	// Rename the binary from mysticeti to hammerhead2
	cmd := exec.Command("mv", "protocols/hammerhead2/assets/mysticeti", "protocols/hammerhead2/assets/hammerhead2")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic("Error renaming binary: " + err.Error() + " " + string(output))
	}

	var wg sync.WaitGroup
	wg.Add(int(num_replicas_int))

	for j := int64(0); j < num_replicas_int; j++ {
		go func(i int) {
			nodes[i].Put_Load("protocols/hammerhead2/assets/hammerhead2", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/hammerhead2/assets/config-rewrite.py", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/hammerhead2/assets/performance_graph.py", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			wg.Done()
		}(int(j))
	}
	wg.Wait()
	ba.logger.Debug(fmt.Sprintf("Copied the hammerhead2 binary to all the nodes\n"), 0)
	return nil
}

func (ba *Hammerhead2) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {

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
			nodes[j].ExecCmd("pkill -KILL -f hammerhead2")
			nodes[j].ExecCmd(fmt.Sprintf("rm -r %vbench/logs/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vbench/logs/", nodes[j].HomeDir))
			wg.Done()
		}(i)
	}
	wg.Wait()

	ba.logger.Debug(fmt.Sprintf("Killed all the replicas\n"), 0)

	int_load, _ := strconv.Atoi(param_load)
	load := strconv.Itoa(int_load / int(num_replicas))

	wave_length := "3"

	number_of_leaders, ok := ba.options.Option["number_of_leaders"]
	if !ok {
		panic("number_of_leaders not found in options")
	}

	max_leaders_per_round, ok := ba.options.Option["max_leaders_per_round"]
	if !ok {
		panic("max_leaders_per_round not found in options")
	}

	interval, ok := ba.options.Option["interval"]
	if !ok {
		panic("interval not found in options")
	}

	enable_pipelining := "true"

	enable_synchronizer := "true"

	transaction_size := param_size

	sshCmd := exec.Command("python3", []string{"protocols/hammerhead2/assets/genrate-configs.py", "--wave_length", wave_length, "--number_of_leaders", number_of_leaders, "--enable_pipelining", enable_pipelining, "--consensus_only", "true", "--enable_synchronizer", enable_synchronizer, "--initial_delay_secs", "5", "--initial_delay_nanos", "0", "--load", load, "--transaction_size", transaction_size, "--max_leaders_per_round", max_leaders_per_round, "--interval", interval, "--output_dir", "protocols/hammerhead2/assets/"}...)
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
			nodes[j].ExecCmd("pkill -KILL -f hammerhead2")
			nodes[j].ExecCmd(fmt.Sprintf("rm %vclient-times-%v.txt", nodes[j].HomeDir, j))
			nodes[j].Put_Load("protocols/hammerhead2/assets/client-parameters.yml", fmt.Sprintf("%vbench/", nodes[j].HomeDir))
			nodes[j].Put_Load("protocols/hammerhead2/assets/node-parameters.yml", fmt.Sprintf("%vbench/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("rm -rf %vbench/storage-{0..%v}", nodes[j].HomeDir, num_replicas-1))
			nodes[j].ExecCmd(fmt.Sprintf("./bench/hammerhead2 benchmark-genesis --ips %v --working-directory %v --node-parameters-path %vnode-parameters.yml", ip_string, nodes[j].HomeDir+"bench/", nodes[j].HomeDir+"bench/"))
			nodes[j].ExecCmd(fmt.Sprintf("python3 %vbench/config-rewrite.py %vbench/public-config.yaml", nodes[j].HomeDir, nodes[j].HomeDir))
			wg1.Done()
		}(i)
	}
	wg1.Wait()

	ba.logger.Debug(fmt.Sprintf("Generated the node private keys"), 0)

	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].ExecCmd(fmt.Sprintf("./bench/hammerhead2 run --authority %v --committee-path %vbench/committee.yaml --public-config-path %vbench/public-config.yaml --private-config-path %vbench/private-config-%v.yaml --client-parameters-path %vbench/client-parameters.yml", j, nodes[j].HomeDir, nodes[j].HomeDir, nodes[j].HomeDir, j, nodes[j].HomeDir))
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
			nodes[j].ExecCmd("pkill -KILL -f hammerhead2")
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
	perf_result := make([]string, num_replicas)
	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			perf_result[j] = nodes[j].ExecCmd(fmt.Sprintf("python3 %vbench/performance_graph.py hammerhead2  %v %vclient-times-%v.txt", nodes[j].HomeDir, duration, nodes[j].HomeDir, j))
			nodes[j].Get_Load(fmt.Sprintf("%vbench/logs/hammerhead2_latency.pdf", nodes[j].HomeDir), fmt.Sprintf("logs/%v_latency.pdf", j))
			nodes[j].Get_Load(fmt.Sprintf("%vbench/logs/hammerhead2_throughput.pdf", nodes[j].HomeDir), fmt.Sprintf("logs/%v_throughput.pdf", j))
			wg3.Done()
		}(i)
	}
	wg3.Wait()

	ba.logger.Debug(fmt.Sprintf("calculated performance"), 0)

	result <- ba.getPerformance(perf_result)

	var wg5 sync.WaitGroup
	wg5.Add(int(num_replicas))
	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].ExecCmd("pkill -KILL -f hammerhead2")
			nodes[j].ExecCmd(fmt.Sprintf("rm %vclient-times-%v.txt", nodes[j].HomeDir, j))
			nodes[j].ExecCmd(fmt.Sprintf("rm -rf %vbench/storage-{0..%v}", nodes[j].HomeDir, num_replicas-1))
			wg5.Done()
		}(i)
	}
	wg5.Wait()

}

func (ba *Hammerhead2) ExtractOptions(path string) protocols.ConsensusOptions {
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

	ba.logger.Debug(fmt.Sprintf("Hammerhead2 options:\n %v\n", options.Option), 0)

	ba.options = options
	return options
}

func (ba *Hammerhead2) getPerformance(outputs []string) util.Performance {
	p := util.Performance{
		Option: make(map[string]string),
	}
	var maxTx float64
	var maxLat float64

	for _, out := range outputs {
		parts := strings.Split(out, " ")
		if len(parts) < 2 {
			continue
		}

		tx, err1 := strconv.ParseFloat(parts[0], 64)
		lat, err2 := strconv.ParseFloat(parts[1], 64)

		if err1 != nil || err2 != nil {
			continue
		}

		if tx > maxTx {
			maxTx = tx
			maxLat = lat
		}
	}

	p.Option["throughput"] = fmt.Sprintf("%.2f requests per second", maxTx)
	p.Option["average latency"] = fmt.Sprintf("%.2f ms", maxLat)

	fmt.Printf("%v,%v,%v,", maxTx, maxLat, 0)

	return p
}
