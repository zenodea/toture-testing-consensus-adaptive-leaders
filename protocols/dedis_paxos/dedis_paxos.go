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

type Dedis_Paxos struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewDedis_Paxos(logger *util.Logger) *Dedis_Paxos {
	return &Dedis_Paxos{
		logger: logger,
	}
}

func (ba *Dedis_Paxos) CopyConsensus(nodes []*common.Node) error {
	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("num_replicas not found in options")
	}
	num_clients, ok := ba.options.Option["num_clients"]
	if !ok {
		panic("num_clients not found in options")
	}

	num_replicas_int, err := strconv.ParseInt(num_replicas, 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}
	num_clients_int, err := strconv.ParseInt(num_clients, 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_clients")
	}

	if num_replicas_int+num_clients_int > int64(len(nodes)) {
		panic("Not enough nodes to deploy dedis_paxos")
	}

	config_inputs := []string{"protocols/dedis_paxos/assets/config-generate.py", num_replicas, num_clients}
	for i := int64(0); i < num_clients_int+num_replicas_int; i++ {
		config_inputs = append(config_inputs, nodes[i].Ip)
	}

	fmt.Printf("Running python command: %v\n", config_inputs)

	sshCmd := exec.Command("python3", config_inputs...)
	output, err := sshCmd.CombinedOutput()
	if err != nil {

		panic("Error while running config-generate.py " + err.Error() + " " + string(output))
	}
	if len(output) > 0 {
		// write the output to protocols/Dedis_Paxos/assets/ip_config.yaml
		err = ioutil.WriteFile("protocols/dedis_paxos/assets/ip_config.yaml", output, 0644)
		if err != nil {
			panic("Error while writing to ip_config.yaml " + err.Error())
		} else {
			fmt.Printf("ip_config.yaml written successfully with content:\n %v\n", string(output))
		}
	}

	// copy the replica binary, client binary and configuration file to the nodes

	var wg sync.WaitGroup
	wg.Add(int(num_clients_int + num_replicas_int))

	for j := int64(0); j < num_clients_int+num_replicas_int; j++ {
		go func(i int) {
			nodes[i].Put_Load("protocols/dedis_paxos/assets/ip_config.yaml", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/dedis_paxos/assets/replica", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/dedis_paxos/assets/client", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			wg.Done()
		}(int(j))
	}
	wg.Wait()
	fmt.Print("Copied the dedis_paxos binaries to all the nodes\n")

	return nil
}

func (ba *Dedis_Paxos) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
	replica_path := "/bench/replica"
	ctl_path := "/bench/client"

	num_replicas, err := strconv.ParseInt(ba.options.Option["num_replicas"], 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}
	num_clients, err := strconv.ParseInt(ba.options.Option["num_clients"], 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_clients")
	}

	view_timeout_time, ok := ba.options.Option["view_timeout_time"]
	if !ok {
		panic(err.Error() + " while parsing view_timeout_time")
	}

	arrival_rate, ok := ba.options.Option["arrival_rate"]
	if !ok {
		panic(err.Error() + " while parsing arrival_rate")
	}

	pipeline_length, ok := ba.options.Option["pipeline_length"]
	if !ok {
		panic(err.Error() + " while parsing pipeline_length")
	}

	var wg sync.WaitGroup
	wg.Add(int(num_replicas + num_clients))
	for i := 0; i < int(num_replicas+num_clients); i++ {
		go func(j int) {
			nodes[j].ExecCmd("pkill -KILL -f replica")
			nodes[j].ExecCmd("pkill -KILL -f client")
			nodes[j].ExecCmd(fmt.Sprintf("rm -r %vbench/logs/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vbench/logs/", nodes[j].HomeDir))
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Print("Killed all the replicas and clients\n")

	for j := 0; j < int(num_replicas); j++ {
		go func(i int) {
			nodes[i].ExecCmd("." + replica_path + " --name " + strconv.Itoa(i+1) + " --consAlgo paxos " + " --viewTimeout " + view_timeout_time + " --pipelineLength " + pipeline_length + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[i].HomeDir))
		}(j)
	}

	time.Sleep(5 * time.Second)

	fmt.Print("Started all the replicas\n")

	nodes[num_replicas].ExecCmd("." + ctl_path + " --name " + strconv.Itoa(51) + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[num_replicas].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[num_replicas].HomeDir) + " --requestType status --operationType 1 ")

	fmt.Print("Sent initial status to bootstrap\n")

	time.Sleep(15 * time.Second)

	nodes[num_replicas].ExecCmd("." + ctl_path + " --name " + strconv.Itoa(51) + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[num_replicas].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[num_replicas].HomeDir) + " --requestType status --operationType 3 ")

	fmt.Print("Sent consensus start-up to bootstrap\n")

	time.Sleep(15 * time.Second)

	clientOutputs := make([]string, num_clients)
	m := 1
	for j := int(num_replicas); j < int(num_replicas+num_clients); j++ {
		go func(i int, k int) {
			clientOutputs[i-int(num_replicas)] = nodes[i].ExecCmd("." + ctl_path + " --name " + strconv.Itoa(50+k) + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[i].HomeDir) + " --requestType request --arrivalRate  " + arrival_rate + " --testDuration " + strconv.Itoa(duration))
		}(j, m)
		m++
	}

	fmt.Print("Started all the clients\n")

	time.Sleep(6 * time.Second)

	fmt.Print("Bootstrap complete\n")
	bootstrap_complete <- true

	time.Sleep(time.Duration(2*duration) * time.Second)

	fmt.Print("Finished the clients\n")

	var wg1 sync.WaitGroup
	wg1.Add(int(num_replicas + num_clients))
	for j := 0; j < int(num_replicas+num_clients); j++ {
		go func(i int) {
			nodes[i].ExecCmd("pkill -KILL -f replica")
			nodes[i].ExecCmd("pkill -KILL -f client")
			wg1.Done()
		}(j)
	}
	wg1.Wait()

	fmt.Print("Killed all the replicas and clients\n")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Error getting home directory:" + err.Error())
	}

	sshCmd := exec.Command("rm", []string{"-r", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err := sshCmd.CombinedOutput()
	if err != nil {
		print("Error while deleting logs/ " + err.Error() + " " + string(output) + "\n")
	} else {
		print("deleted local logs/ successfully\n" + string(output) + "\n")
	}

	sshCmd = exec.Command("mkdir", []string{"-p", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		panic("Error while creating logs/ " + err.Error() + " " + string(output) + "\n")
	} else {
		print("created logs/ successfully\n" + string(output) + "\n")
	}

	var wg2 sync.WaitGroup
	wg2.Add(int(num_clients))
	m = 1
	for j := int(num_replicas); j < int(num_replicas+num_clients); j++ {
		go func(i int, k int) {
			nodes[i].Get_Load(fmt.Sprintf("%vbench/logs/%v.txt", nodes[i].HomeDir, 50+k), "logs/")
			wg2.Done()
		}(j, m)
		m++
	}
	wg2.Wait()
	fmt.Println("Downloaded all the dedis_paxos client logs")

	command := "protocols/dedis_paxos/assets/performance_graph.py"
	file_names := []string{}

	m = 1
	for j := int(num_replicas); j < int(num_replicas+num_clients); j++ {
		logFile := filepath.Join(homeDir, fmt.Sprintf("toture-testing-consensus/logs/%v.txt", 50+m))

		file_names = append(file_names, logFile)

		sshCmd = exec.Command("python3", []string{command, "dedis_paxos-" + strconv.Itoa(50+m), logFile}...)
		output, err = sshCmd.CombinedOutput()
		if err != nil {
			print("Error while generating performance graph " + err.Error() + " " + string(output) + "\n")
		} else {
			print("Dedis_Paxos Performance graph generated successfully\n" + string(output) + "\n")
		}

		m++
	}

	sshCmd = exec.Command("python3", append([]string{command, "dedis_paxos"}, file_names...)...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		print("Error while generating performance graphs " + err.Error() + " " + string(output) + "\n")
	} else {
		print("Dedis_Paxos Performance graphs generated successfully\n" + string(output) + "\n")
	}

	result <- ba.GetPerformance(clientOutputs)
}

func (ba *Dedis_Paxos) ExtractOptions(path string) protocols.ConsensusOptions {
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

	fmt.Printf("dedis_paxos options:\n %v\n", options.Option)

	ba.options = options
	return options
}

func (ba *Dedis_Paxos) GetPerformance(outputs []string) util.Performance {
	throughput := make([]float64, len(outputs))
	medians := make([]float64, len(outputs))
	percentil99s := make([]float64, len(outputs))

	for i := 0; i < len(outputs); i++ {
		outputLines := strings.Split(outputs[i], "\n")
		for j := 0; j < len(outputLines); j++ {
			if strings.Contains(outputLines[j], "Throughput") {
				throughput[i], _ = strconv.ParseFloat(strings.Split(outputLines[j], " ")[5], 64)
			}
			if strings.Contains(outputLines[j], "Median") {
				medians[i], _ = strconv.ParseFloat(strings.Split(outputLines[j], " ")[3], 64)
			}
			if strings.Contains(outputLines[j], "99 pecentile") {
				percentil99s[i], _ = strconv.ParseFloat(strings.Split(outputLines[j], " ")[4], 64)
			}
		}
	}
	sum_throughput := 0.0
	sum_median := 0.0
	sum_percentle := 0.0

	for i := 0; i < len(throughput); i++ {
		sum_throughput += throughput[i]
		sum_median += medians[i]
		sum_percentle += percentil99s[i]
	}

	return util.Performance{
		map[string]string{
			"throughput":   fmt.Sprintf("%v", sum_throughput),
			"median":       fmt.Sprintf("%v", sum_median/float64(len(throughput))),
			"percentile99": fmt.Sprintf("%v", sum_percentle/float64(len(throughput))),
		},
	}
}
