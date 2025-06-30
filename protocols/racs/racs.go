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

type Racs struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewRacs(logger *util.Logger) *Racs {
	return &Racs{
		logger: logger,
	}
}

func (ba *Racs) CopyConsensus(nodes []*common.Node) error {
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
		panic("Not enough nodes to deploy racs")
	}

	config_inputs := []string{"protocols/racs/assets/config-generate.py", num_replicas, num_clients}
	for i := int64(0); i < num_clients_int+num_replicas_int; i++ {
		config_inputs = append(config_inputs, nodes[i].Ip)
	}

	ba.logger.Debug(fmt.Sprintf("Running python command: %v\n", config_inputs), 0)

	sshCmd := exec.Command("python3", config_inputs...)
	output, err := sshCmd.CombinedOutput()
	if err != nil {

		panic("Error while running config-generate.py " + err.Error() + " " + string(output))
	}
	if len(output) > 0 {
		// write the output to protocols/Racs/assets/ip_config.yaml
		err = ioutil.WriteFile("protocols/racs/assets/ip_config.yaml", output, 0644)
		if err != nil {
			panic("Error while writing to ip_config.yaml " + err.Error())
		} else {
			ba.logger.Debug(fmt.Sprintf("ip_config.yaml written successfully with content:\n %v\n", string(output)), 0)
		}
	}

	// copy the replica binary, client binary and configuration file to the nodes

	var wg sync.WaitGroup
	wg.Add(int(num_clients_int + num_replicas_int))

	for j := int64(0); j < num_clients_int+num_replicas_int; j++ {
		go func(i int) {
			nodes[i].Put_Load("protocols/racs/assets/ip_config.yaml", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/racs/assets/replica", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/racs/assets/client", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			wg.Done()
		}(int(j))
	}
	wg.Wait()
	ba.logger.Debug(fmt.Sprintf("Copied the racs binaries to all the nodes\n"), 0)

	return nil
}

func (ba *Racs) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
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

	pipeline_length, ok := ba.options.Option["pipeline_length"]
	if !ok {
		panic(err.Error() + " while parsing pipeline_length")
	}

	arrival_rate, ok := ba.options.Option["arrival_rate"]
	if !ok {
		panic(err.Error() + " while parsing arrival_rate")
	}

	network_batch_time, ok := ba.options.Option["network_batch_time"]
	if !ok {
		panic(err.Error() + " while parsing network_batch_time")
	}

	replica_batch_size, ok := ba.options.Option["replica_batch_size"]
	if !ok {
		panic(err.Error() + " while parsing replica_batch_size")
	}

	replica_batch_time, ok := ba.options.Option["replica_batch_time"]
	if !ok {
		panic(err.Error() + " while parsing replica_batch_time")
	}

	key_len, ok := ba.options.Option["key_len"]
	if !ok {
		panic(err.Error() + " while parsing key_len")
	}

	val_len, ok := ba.options.Option["val_len"]
	if !ok {
		panic(err.Error() + " while parsing val_len")
	}

	client_batch_size, ok := ba.options.Option["client_batch_size"]
	if !ok {
		panic(err.Error() + " while parsing client_batch_size")
	}

	client_batch_time, ok := ba.options.Option["client_batch_time"]
	if !ok {
		panic(err.Error() + " while parsing client_batch_time")
	}

	client_window, ok := ba.options.Option["client_window"]
	if !ok {
		panic(err.Error() + " while parsing client_window")
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

	ba.logger.Debug(fmt.Sprintf("Killed all the replicas and clients\n"), 0)

	for j := 0; j < int(num_replicas); j++ {
		go func(i int) {
			nodes[i].ExecCmd("." + replica_path + " --name " + strconv.Itoa(i+1) + " --viewTimeout " + view_timeout_time + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[i].HomeDir) + " --pipelineLength " + pipeline_length + " --networkbatchTime " + network_batch_time + " --batchSize " + replica_batch_size + " --batchTime  " + replica_batch_time + " --keyLen " + key_len + " --valLen  " + val_len)
		}(j)
	}

	time.Sleep(5 * time.Second)

	ba.logger.Debug(fmt.Sprintf("Started all the replicas\n"), 0)

	nodes[num_replicas].ExecCmd("." + ctl_path + " --name " + strconv.Itoa(51) + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[num_replicas].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[num_replicas].HomeDir) + " --requestType status --operationType 1 ")

	ba.logger.Debug(fmt.Sprintf("Sent initial status to bootstrap\n"), 0)

	time.Sleep(15 * time.Second)

	nodes[num_replicas].ExecCmd("." + ctl_path + " --name " + strconv.Itoa(51) + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[num_replicas].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[num_replicas].HomeDir) + " --requestType status --operationType 3 ")

	ba.logger.Debug(fmt.Sprintf("Sent consensus start-up to bootstrap\n"), 0)

	time.Sleep(15 * time.Second)

	clientOutputs := make([]string, num_clients)
	m := 1
	for j := int(num_replicas); j < int(num_replicas+num_clients); j++ {
		go func(i int, k int) {
			clientOutputs[i-int(num_replicas)] = nodes[i].ExecCmd("." + ctl_path + " --name " + strconv.Itoa(50+k) + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[i].HomeDir) + " --requestType request --arrivalRate  " + arrival_rate + " --testDuration " + strconv.Itoa(duration) + " --keyLen " + key_len + " --valLen " + val_len + " --batchSize " + client_batch_size + " --batchTime " + client_batch_time + " --window " + client_window)
		}(j, m)
		m++
	}

	ba.logger.Debug(fmt.Sprintf("Started all the clients\n"), 0)

	time.Sleep(10 * time.Second)

	ba.logger.Debug(fmt.Sprintf("Bootstrap complete\n"), 0)
	bootstrap_complete <- true

	time.Sleep(time.Duration(3*duration) * time.Second)

	ba.logger.Debug(fmt.Sprintf("Finished the clients\n"), 0)

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
	ba.logger.Debug(fmt.Sprintf("Downloaded all the racs client logs"), 0)

	command := "protocols/racs/assets/performance_graph.py"
	file_names := []string{}

	m = 1
	for j := int(num_replicas); j < int(num_replicas+num_clients); j++ {
		logFile := filepath.Join(homeDir, fmt.Sprintf("toture-testing-consensus/logs/%v.txt", 50+m))

		file_names = append(file_names, logFile)

		sshCmd = exec.Command("python3", []string{command, "racs-" + strconv.Itoa(50+m), logFile}...)
		output, err = sshCmd.CombinedOutput()
		if err != nil {
			ba.logger.Debug(fmt.Sprintf("Error while generating performance graph "+err.Error()+" "+string(output)+"\n"), 0)
		} else {
			ba.logger.Debug(fmt.Sprintf("Racs Performance graph generated successfully\n"+string(output)+"\n"), 0)
		}

		m++
	}

	sshCmd = exec.Command("python3", append([]string{command, "racs"}, file_names...)...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug(fmt.Sprintf("Error while generating performance graphs "+err.Error()+" "+string(output)+"\n"), 0)
	} else {
		ba.logger.Debug(fmt.Sprintf("Racs Performance graphs generated successfully\n"+string(output)+"\n"), 0)
	}

	result <- ba.GetPerformance(clientOutputs)
}

func (ba *Racs) ExtractOptions(path string) protocols.ConsensusOptions {
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

	ba.logger.Debug(fmt.Sprintf("racs options:\n %v\n", options.Option), 0)

	ba.options = options
	return options
}

func (ba *Racs) GetPerformance(outputs []string) util.Performance {
	throughput := make([]float64, len(outputs))
	medians := make([]float64, len(outputs))
	percentil99s := make([]float64, len(outputs))

	for i := 0; i < len(outputs); i++ {
		outputLines := strings.Split(outputs[i], "\n")
		for j := 0; j < len(outputLines); j++ {
			if strings.Contains(outputLines[j], "Throughput") {
				throughput[i], _ = strconv.ParseFloat(strings.Split(outputLines[j], " ")[2], 64)
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

	fmt.Printf("%v ", p)

	return p
}
