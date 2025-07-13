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

type Efficient struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewEfficient(logger *util.Logger) *Efficient {
	return &Efficient{
		logger: logger,
	}
}

func (ba *Efficient) CopyConsensus(nodes []*common.Node) error {
	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("num_replicas not found in options")
	}

	num_replicas_int, err := strconv.ParseInt(num_replicas, 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}

	if num_replicas_int > int64(len(nodes)) {
		panic("Not enough nodes to deploy efficient")
	}

	// copy the replica binary, client binary

	var wg sync.WaitGroup
	wg.Add(int(num_replicas_int))

	for j := int64(0); j < num_replicas_int; j++ {
		go func(i int) {
			nodes[i].Put_Load("protocols/efficient/assets/epaxos_client", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/efficient/assets/epaxos_master", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/efficient/assets/epaxos_server", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			wg.Done()
		}(int(j))
	}
	wg.Wait()
	ba.logger.Debug(fmt.Sprintf("Copied the efficient binaries to all the nodes\n"), 0)

	return nil
}

func (ba *Efficient) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
	data, err := os.ReadFile("params.param")
	if err != nil {
		panic(err.Error())
	}

	parts := strings.Fields(string(data))
	if len(parts) != 2 {
		panic("Expected exactly 2 values in params.param")
	}

	param_load := parts[0]

	_ = parts[1]

	replica_path := "/bench/epaxos_server"
	ctl_path := "/bench/epaxos_client"
	master_path := "/bench/epaxos_master"

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

	algo, ok := ba.options.Option["algo"]
	if !ok {
		panic(err.Error() + " while parsing algo")
	}

	execF, ok := ba.options.Option["exec"]
	if !ok {
		panic(err.Error() + " while parsing exec")
	}

	dreply, ok := ba.options.Option["dreply"]
	if !ok {
		panic(err.Error() + " while parsing dreply")
	}

	durable, ok := ba.options.Option["durable"]
	if !ok {
		panic(err.Error() + " while parsing durable")
	}

	thrifty, ok := ba.options.Option["thrifty"]
	if !ok {
		panic(err.Error() + " while parsing thrifty")
	}

	writes, ok := ba.options.Option["w"]
	if !ok {
		panic(err.Error() + " while parsing w")
	}

	conflicts, ok := ba.options.Option["c"]
	if !ok {
		panic(err.Error() + " while parsing c")
	}

	var wg sync.WaitGroup
	wg.Add(int(num_replicas))
	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].ExecCmd("pkill -KILL -f epaxos_client")
			nodes[j].ExecCmd("pkill -KILL -f epaxos_master")
			nodes[j].ExecCmd("pkill -KILL -f epaxos_server")
			nodes[j].ExecCmd(fmt.Sprintf("rm -r %vbench/logs/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vbench/logs/", nodes[j].HomeDir))
			wg.Done()
		}(i)
	}
	wg.Wait()

	ba.logger.Debug(fmt.Sprintf("Killed all the replicas and clients\n"), 0)

	go nodes[0].ExecCmd("." + master_path + " -N " + strconv.Itoa(int(num_replicas)))

	time.Sleep(5 * time.Second)

	ba.logger.Debug(fmt.Sprintf("Started the master\n"), 0)

	for i := 0; i < int(num_replicas); i++ {
		go nodes[i].ExecCmd("." + replica_path + " -port 10000 " + " -maddr " + nodes[0].Ip + " -addr " + nodes[i].Ip + " -batchSize 3000 " + " -batchTime 5000 " + " -pipeline " + pipeline_length + " " + algo + " " + execF + " " + dreply + " " + durable + " " + thrifty)
		time.Sleep(5 * time.Second)
	}

	ba.logger.Debug(fmt.Sprintf("Started all the replicas\n"), 0)

	time.Sleep(5 * time.Second)

	clientOutputs := make([]string, num_replicas)
	m := 1
	for j := 0; j < int(num_replicas); j++ {
		go func(i int, k int) {
			if algo != "-pa" {
				clientOutputs[i] = nodes[i].ExecCmd("." + ctl_path + " -name " + strconv.Itoa(50+k) + " -maddr " + nodes[0].Ip + "  -clientBatchSize 50  -logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " -arrivalRate  " + arrival_rate + " -testDuration " + strconv.Itoa(duration) + " -leaderTimeout " + view_timeout_time + " -defaultReplica " + strconv.Itoa(i) + " -w " + writes + " -c " + conflicts)
			} else {
				clientOutputs[i] = nodes[i].ExecCmd("." + ctl_path + " -name " + strconv.Itoa(50+k) + " -maddr " + nodes[0].Ip + "  -clientBatchSize 50  -logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " -arrivalRate  " + arrival_rate + " -testDuration " + strconv.Itoa(duration) + " -leaderTimeout " + view_timeout_time + " -defaultReplica " + strconv.Itoa(i) + " -w " + writes + " -c " + conflicts + " -l")
			}
		}(j, m)
		m++
	}

	ba.logger.Debug(fmt.Sprintf("Started all the clients\n"), 0)

	time.Sleep(5 * time.Second)

	ba.logger.Debug(fmt.Sprintf("Bootstrap complete\n"), 0)
	bootstrap_complete <- true

	time.Sleep(time.Duration(2*duration) * time.Second)

	ba.logger.Debug(fmt.Sprintf("Finished the clients\n"), 0)

	var wg1 sync.WaitGroup
	wg1.Add(int(num_replicas))
	for j := 0; j < int(num_replicas); j++ {
		go func(i int) {
			nodes[i].ExecCmd("pkill -KILL -f epaxos_client")
			nodes[i].ExecCmd("pkill -KILL -f epaxos_master")
			nodes[i].ExecCmd("pkill -KILL -f epaxos_server")
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
	wg2.Add(int(num_replicas))
	m = 1
	for j := 0; j < int(num_replicas); j++ {
		go func(i int, k int) {
			nodes[i].Get_Load(fmt.Sprintf("%vbench/logs/%v.txt", nodes[i].HomeDir, 50+k), "logs/")
			wg2.Done()
		}(j, m)
		m++
	}
	wg2.Wait()
	ba.logger.Debug(fmt.Sprintf("Downloaded all the efficient client logs"), 0)

	command := "protocols/efficient/assets/performance_graph.py"
	file_names := []string{}

	m = 1
	for j := 0; j < int(num_replicas); j++ {
		logFile := filepath.Join(homeDir, fmt.Sprintf("toture-testing-consensus/logs/%v.txt", 50+m))
		file_names = append(file_names, logFile)
		m++
	}

	sshCmd = exec.Command("python3", append([]string{command, "efficient"}, file_names...)...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug(fmt.Sprintf("Error while generating performance graphs "+err.Error()+" "+string(output)+"\n"), 0)
	} else {
		ba.logger.Debug(fmt.Sprintf("Efficient Performance graphs generated successfully\n"+string(output)+"\n"), 0)
	}

	result <- ba.GetPerformance(clientOutputs)
}

func (ba *Efficient) ExtractOptions(path string) protocols.ConsensusOptions {
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

	_, ok := options.Option["exec"]
	if !ok {
		options.Option["exec"] = ""
	}

	_, ok = options.Option["dreply"]
	if !ok {
		options.Option["dreply"] = ""
	}

	_, ok = options.Option["durable"]
	if !ok {
		options.Option["durable"] = ""
	}

	_, ok = options.Option["thrifty"]
	if !ok {
		options.Option["thrifty"] = ""
	}

	ba.logger.Debug(fmt.Sprintf("efficient options:\n %v\n", options.Option), 0)

	ba.options = options
	return options
}

func (ba *Efficient) GetPerformance(outputs []string) util.Performance {
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
