package consensus

import (
	"fmt"
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

	"gopkg.in/yaml.v2"
)

type Dedis_Raft struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewDedis_Raft(logger *util.Logger) *Dedis_Raft {
	return &Dedis_Raft{
		logger: logger,
	}
}

func (ba *Dedis_Raft) CopyConsensus(nodes []*common.Node) error {
	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("num_replicas not found in options")
	}

	num_replicas_int, err := strconv.ParseInt(num_replicas, 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}

	if num_replicas_int > int64(len(nodes)) {
		panic("Not enough nodes to deploy dedis_raft")
	}

	config_inputs := []string{"protocols/dedis_raft/assets/config-generate.py", num_replicas, num_replicas}
	for i := int64(0); i < num_replicas_int; i++ {
		config_inputs = append(config_inputs, nodes[i].Ip)
	}
	for i := int64(0); i < num_replicas_int; i++ {
		config_inputs = append(config_inputs, nodes[i].Ip)
	}

	ba.logger.Debug(fmt.Sprintf("Running python command: %v\n", config_inputs), 0)

	sshCmd := exec.Command("python3", config_inputs...)
	output, err := sshCmd.CombinedOutput()
	if err != nil {

		panic("Error while running config-generate.py " + err.Error() + " " + string(output))
	}
	if len(output) > 0 {
		// write the output to protocols/Dedis_Raft/assets/ip_config.yaml
		err = ioutil.WriteFile("protocols/dedis_raft/assets/ip_config.yaml", output, 0644)
		if err != nil {
			panic("Error while writing to ip_config.yaml " + err.Error())
		} else {
			ba.logger.Debug(fmt.Sprintf("ip_config.yaml written successfully with content:\n %v\n", string(output)), 0)
		}
	}

	// copy the replica binary, client binary and configuration file to the nodes

	var wg sync.WaitGroup
	wg.Add(int(num_replicas_int))

	for j := int64(0); j < num_replicas_int; j++ {
		go func(i int) {
			nodes[i].Put_Load("protocols/dedis_raft/assets/ip_config.yaml", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/dedis_raft/assets/replica", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/dedis_raft/assets/client", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			wg.Done()
		}(int(j))
	}
	wg.Wait()
	ba.logger.Debug(fmt.Sprintf("Copied the dedis_raft binaries to all the nodes\n"), 0)

	return nil
}

func (ba *Dedis_Raft) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {

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

	replica_path := "/bench/replica"
	ctl_path := "/bench/client"

	num_replicas, err := strconv.ParseInt(ba.options.Option["num_replicas"], 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}

	view_timeout_time, ok := ba.options.Option["view_timeout_time"]
	if !ok {
		panic(err.Error() + " while parsing view_timeout_time")
	}

	int_load, _ := strconv.Atoi(param_load)
	arrival_rate := strconv.Itoa(int_load / int(num_replicas))

	pipeline_length, ok := ba.options.Option["pipeline_length"]
	if !ok {
		panic(err.Error() + " while parsing pipeline_length")
	}

	replica_batch_size, ok := ba.options.Option["replica_batch_size"]
	if !ok {
		panic("replica_batch_size not found in options")
	}

	replica_batch_time, ok := ba.options.Option["replica_batch_time"]
	if !ok {
		panic("replica_batch_time not found in options")
	}

	int_size, _ := strconv.Atoi(param_size)

	key_len := strconv.Itoa(int_size / 2)

	val_len := strconv.Itoa(int_size / 2)

	client_batch_size, ok := ba.options.Option["client_batch_size"]
	if !ok {
		panic("client_batch_size not found in options")
	}

	client_batch_time, ok := ba.options.Option["client_batch_time"]
	if !ok {
		panic("client_batch_time not found in options")
	}

	client_window, ok := ba.options.Option["client_window"]
	if !ok {
		panic("client_window not found in options")
	}

	var wg sync.WaitGroup
	wg.Add(int(num_replicas))
	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].ExecCmd("pkill -KILL -f replica")
			nodes[j].ExecCmd("pkill -KILL -f client")
			nodes[j].ExecCmd(fmt.Sprintf("rm -r %vbench/logs/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vbench/logs/", nodes[j].HomeDir))
			wg.Done()
		}(i)
	}
	wg.Wait()

	ba.logger.Debug(fmt.Sprintf("Killed all the replicas and clients\n"), 0)

	for j := 0; j < int(num_replicas); j++ {
		go func(i int) {
			nodes[i].ExecCmd("." + replica_path + " --name " + strconv.Itoa(i+1) + " --consAlgo raft " + " --viewTimeout " + view_timeout_time + " --pipelineLength " + pipeline_length + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[i].HomeDir) + "  --batchSize " + replica_batch_size + "  --batchTime " + replica_batch_time + "  --keyLen " + key_len + "  --valLen " + val_len)
		}(j)
	}

	time.Sleep(5 * time.Second)

	ba.logger.Debug(fmt.Sprintf("Started all the replicas\n"), 0)

	nodes[0].ExecCmd("." + ctl_path + " --name " + strconv.Itoa(51) + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[0].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[0].HomeDir) + " --requestType status --operationType 1 ")

	ba.logger.Debug(fmt.Sprintf("Sent initial status to bootstrap\n"), 0)

	time.Sleep(15 * time.Second)

	nodes[0].ExecCmd("." + ctl_path + " --name " + strconv.Itoa(51) + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[0].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[0].HomeDir) + " --requestType status --operationType 3 ")

	ba.logger.Debug(fmt.Sprintf("Sent consensus start-up to bootstrap\n"), 0)

	time.Sleep(15 * time.Second)

	clientOutputs := make([]string, num_replicas)
	m := 1
	for j := 0; j < int(num_replicas); j++ {
		go func(i int, k int) {
			clientOutputs[i] = nodes[i].ExecCmd("." + ctl_path + " --name " + strconv.Itoa(50+k) + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[i].HomeDir) + " --requestType request --arrivalRate  " + arrival_rate + " --testDuration " + strconv.Itoa(duration) + "  --batchSize " + client_batch_size + "  --batchTime " + client_batch_time + "  --keyLen " + key_len + "  --valLen " + val_len + "  --window " + client_window)
		}(j, m)
		m++
	}

	ba.logger.Debug(fmt.Sprintf("Started all the clients\n"), 0)

	time.Sleep(10 * time.Second)

	ba.logger.Debug(fmt.Sprintf("Bootstrap complete\n"), 0)
	bootstrap_complete <- true

	time.Sleep(time.Duration(2*duration) * time.Second)

	ba.logger.Debug(fmt.Sprintf("Finished the clients\n"), 0)

	var wg1 sync.WaitGroup
	wg1.Add(int(num_replicas))
	for j := 0; j < int(num_replicas); j++ {
		go func(i int) {
			nodes[i].ExecCmd("pkill -KILL -f replica")
			nodes[i].ExecCmd("pkill -KILL -f client")
			wg1.Done()
		}(j)
	}
	wg1.Wait()

	ba.logger.Debug(fmt.Sprintf("Killed all the replicas and clients\n"), 0)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Error getting home directory:" + err.Error())
	}

	sshCmd := exec.Command("rm", []string{"-r", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err := sshCmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug(fmt.Sprintf("Error while deleting logs/ "+err.Error()+" "+string(output)+"\n"), 0)
	} else {
		ba.logger.Debug(fmt.Sprintf("deleted local logs/ successfully\n"+string(output)+"\n"), 0)
	}

	sshCmd = exec.Command("mkdir", []string{"-p", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		panic("Error while creating logs/ " + err.Error() + " " + string(output) + "\n")
	} else {
		ba.logger.Debug(fmt.Sprintf("created logs/ successfully\n"+string(output)+"\n"), 0)
	}

	//var wg2 sync.WaitGroup
	//wg2.Add(int(num_replicas))
	//m = 1
	//for j := 0; j < int(num_replicas); j++ {
	//	go func(i int, k int) {
	//		nodes[i].Get_Load(fmt.Sprintf("%vbench/logs/%v.txt", nodes[i].HomeDir, 50+k), "logs/")
	//		wg2.Done()
	//	}(j, m)
	//	m++
	//}
	//wg2.Wait()
	//ba.logger.Debug(fmt.Sprintf("Downloaded all the dedis_raft client logs"), 0)
	//
	//command := "protocols/dedis_raft/assets/performance_graph.py"
	//file_names := []string{}
	//
	//m = 1
	//for j := 0; j < int(num_replicas); j++ {
	//	logFile := filepath.Join(homeDir, fmt.Sprintf("toture-testing-consensus/logs/%v.txt", 50+m))
	//	file_names = append(file_names, logFile)
	//	m++
	//}
	//
	//sshCmd = exec.Command("python3", append([]string{command, "dedis_raft"}, file_names...)...)
	//output, err = sshCmd.CombinedOutput()
	//if err != nil {
	//	ba.logger.Debug(fmt.Sprintf("Error while generating performance graphs "+err.Error()+" "+string(output)+"\n"), 0)
	//} else {
	//	ba.logger.Debug(fmt.Sprintf("Dedis_Raft Performance graphs generated successfully\n"+string(output)+"\n"), 0)
	//}

	result <- ba.GetPerformance(clientOutputs)
}

func (ba *Dedis_Raft) ExtractOptions(path string) protocols.ConsensusOptions {
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

	ba.logger.Debug(fmt.Sprintf("dedis_raft options:\n %v\n", options.Option), 0)

	ba.options = options
	return options
}

func (ba *Dedis_Raft) GetPerformance(outputs []string) util.Performance {
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

	p := util.Performance{
		map[string]string{
			"throughput":   fmt.Sprintf("%v", sum_throughput),
			"median":       fmt.Sprintf("%v", sum_median/float64(len(throughput))),
			"percentile99": fmt.Sprintf("%v", sum_percentle/float64(len(throughput))),
		},
	}
	fmt.Printf("%v,%v,%v,", sum_throughput, sum_median/float64(1000*len(throughput)), sum_percentle/float64(1000*len(throughput)))
	return p
}
