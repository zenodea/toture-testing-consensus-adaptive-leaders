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
		panic("Not enough nodes to deploy efficient")
	}

	// copy the replica binary, client binary

	var wg sync.WaitGroup
	wg.Add(int(num_clients_int + num_replicas_int))

	for j := int64(0); j < num_clients_int+num_replicas_int; j++ {
		go func(i int) {
			nodes[i].Put_Load("protocols/efficient/assets/epaxos_client", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/efficient/assets/epaxos_master", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/efficient/assets/epaxos_server", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			wg.Done()
		}(int(j))
	}
	wg.Wait()
	fmt.Print("Copied the efficient binaries to all the nodes\n")

	return nil
}

func (ba *Efficient) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
	replica_path := "/bench/epaxos_server"
	ctl_path := "/bench/epaxos_client"
	master_path := "/bench/epaxos_master"

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

	algo, ok := ba.options.Option["algo"]
	if !ok {
		panic(err.Error() + " while parsing algo")
	}

	var wg sync.WaitGroup
	wg.Add(int(num_replicas + num_clients))
	for i := 0; i < int(num_replicas+num_clients); i++ {
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

	fmt.Print("Killed all the replicas and clients\n")

	go nodes[0].ExecCmd("." + master_path + " -N " + strconv.Itoa(int(num_replicas)))

	time.Sleep(5 * time.Second)

	fmt.Print("Started the master\n")

	for i := 0; i < int(num_replicas); i++ {
		go nodes[i].ExecCmd("." + replica_path + " -port 10000 " + " -maddr " + nodes[0].Ip + " -addr " + nodes[i].Ip + " -batchSize 3000 " + " -batchTime 5000 " + " -pipeline " + pipeline_length + " -exec  -dreply " + algo)
		time.Sleep(5 * time.Second)
	}

	fmt.Print("Started all the replicas\n")

	time.Sleep(5 * time.Second)

	clientOutputs := make([]string, num_clients)
	m := 1
	for j := int(num_replicas); j < int(num_replicas+num_clients); j++ {
		go func(i int, k int) {
			if algo != "-pa" {
				clientOutputs[i-int(num_replicas)] = nodes[i].ExecCmd("." + ctl_path + " -name " + strconv.Itoa(50+k) + " -maddr " + nodes[0].Ip + " -w 50 -c 2 -clientBatchSize 50  -logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " -arrivalRate  " + arrival_rate + " -testDuration " + strconv.Itoa(duration) + " -leaderTimeout " + view_timeout_time + " -defaultReplica " + strconv.Itoa(i-int(num_replicas)))
			} else {
				clientOutputs[i-int(num_replicas)] = nodes[i].ExecCmd("." + ctl_path + " -name " + strconv.Itoa(50+k) + " -maddr " + nodes[0].Ip + " -w 50 -c 2 -clientBatchSize 50  -logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " -arrivalRate  " + arrival_rate + " -testDuration " + strconv.Itoa(duration) + " -leaderTimeout " + view_timeout_time + " -defaultReplica " + strconv.Itoa(i-int(num_replicas)) + " -l")
			}
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
			nodes[i].ExecCmd("pkill -KILL -f epaxos_client")
			nodes[i].ExecCmd("pkill -KILL -f epaxos_master")
			nodes[i].ExecCmd("pkill -KILL -f epaxos_server")
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
	fmt.Println("Downloaded all the efficient client logs")

	command := "protocols/efficient/assets/performance_graph.py"
	file_names := []string{}

	m = 1
	for j := int(num_replicas); j < int(num_replicas+num_clients); j++ {
		logFile := filepath.Join(homeDir, fmt.Sprintf("toture-testing-consensus/logs/%v.txt", 50+m))

		file_names = append(file_names, logFile)

		sshCmd = exec.Command("python3", []string{command, "efficient-" + strconv.Itoa(50+m), logFile}...)
		output, err = sshCmd.CombinedOutput()
		if err != nil {
			print("Error while generating performance graph " + err.Error() + " " + string(output) + "\n")
		} else {
			print("Efficient Performance graph generated successfully\n" + string(output) + "\n")
		}
		m++
	}

	sshCmd = exec.Command("python3", append([]string{command, "efficient"}, file_names...)...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		print("Error while generating performance graphs " + err.Error() + " " + string(output) + "\n")
	} else {
		print("Efficient Performance graphs generated successfully\n" + string(output) + "\n")
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

	fmt.Printf("efficient options:\n %v\n", options.Option)

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

	return util.Performance{
		map[string]string{
			"throughput":   fmt.Sprintf("%v", sum_throughput),
			"median":       fmt.Sprintf("%v", sum_median/float64(len(throughput))),
			"percentile99": fmt.Sprintf("%v", sum_percentle/float64(len(throughput))),
		},
	}
}
