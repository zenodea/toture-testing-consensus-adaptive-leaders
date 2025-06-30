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

type Rabia struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewRabia(logger *util.Logger) *Rabia {
	return &Rabia{
		logger: logger,
	}
}

func (ba *Rabia) CopyConsensus(nodes []*common.Node) error {
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
		panic("Not enough nodes to deploy rabia")
	}

	// copy the replica binary

	var wg sync.WaitGroup
	wg.Add(int(num_clients_int + num_replicas_int))

	for j := int64(0); j < num_clients_int+num_replicas_int; j++ {
		go func(i int) {
			nodes[i].Put_Load("protocols/rabia/assets/rabia", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			wg.Done()
		}(int(j))
	}
	wg.Wait()
	ba.logger.Debug(fmt.Sprintf("Copied the rabia binaries to all the nodes\n"), 0)

	return nil
}

func (ba *Rabia) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
	rabia_path := "/bench/rabia"

	num_replicas, err := strconv.ParseInt(ba.options.Option["num_replicas"], 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}
	num_clients, err := strconv.ParseInt(ba.options.Option["num_clients"], 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_clients")
	}

	arrival_rate, ok := ba.options.Option["arrival_rate"]
	if !ok {
		panic(err.Error() + " while parsing arrival_rate")
	}

	rabia_client_batch_size, ok := ba.options.Option["rabia_client_batch_size"]
	if !ok {
		panic(err.Error() + " while parsing rabia_client_batch_size")
	}

	rabia_proxy_batch_size, ok := ba.options.Option["rabia_proxy_batch_size"]
	if !ok {
		panic(err.Error() + " while parsing rabia_proxy_batch_size")
	}

	rabia_proxy_batch_timeout, ok := ba.options.Option["rabia_proxy_batch_timeout"]
	if !ok {
		panic(err.Error() + " while parsing rabia_proxy_batch_timeout")
	}

	var wg sync.WaitGroup
	wg.Add(int(num_replicas + num_clients))
	for i := 0; i < int(num_replicas+num_clients); i++ {
		go func(j int) {
			nodes[j].ExecCmd("pkill -KILL -f rabia")
			nodes[j].ExecCmd(fmt.Sprintf("rm -r %vbench/logs/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vbench/logs/", nodes[j].HomeDir))
			wg.Done()
		}(i)
	}
	wg.Wait()

	ba.logger.Debug(fmt.Sprintf("Killed all the replicas and clients\n"), 0)

	Controller := nodes[0].Ip + ":9000"
	NServers := num_replicas
	NClients := num_clients
	RC_Peers_N := ""

	for i := 0; i < int(num_replicas); i++ {
		RC_Peers_N += nodes[i].Ip + ":10000,"
	}
	RC_Peers_N = RC_Peers_N[0 : len(RC_Peers_N)-1]

	for j := 0; j < int(num_replicas); j++ {
		go func(i int) {
			export_command := fmt.Sprintf(
				"export LogFilePath=%vbench/logs/ RC_Ctrl=%v RC_Folder=%v/bench/ RC_LLevel=\"warn\" Rabia_ClosedLoop=false Rabia_NServers=%v Rabia_NFaulty=0 Rabia_NClients=%v Rabia_NConcurrency=1 Rabia_ClientBatchSize=%v Rabia_ClientTimeout=%v Rabia_ClientThinkTime=0 Rabia_ClientNRequests=0 Rabia_ClientArrivalRate=%v Rabia_ProxyBatchSize=%v Rabia_ProxyBatchTimeout=%v Rabia_NetworkBatchSize=0 Rabia_NetworkBatchTimeout=0 RC_Peers=%v Rabia_StorageMode=0",
				nodes[i].HomeDir, Controller, nodes[i].HomeDir, NServers, NClients, rabia_client_batch_size, duration, arrival_rate, rabia_proxy_batch_size, rabia_proxy_batch_timeout, RC_Peers_N)
			svr_export := fmt.Sprintf("export RC_Role=svr RC_Index=%v RC_SvrIp=\"%v\" RC_PPort=\"11000\" RC_NPort=\"10000\"", i, nodes[i].Ip)
			nodes[i].ExecCmd(svr_export + ";" + export_command + ";" + "." + rabia_path)
			if i == 0 {
				ba.logger.Debug(fmt.Sprintf("export_command: %v\n\n\n", export_command), 0)
				ba.logger.Debug(fmt.Sprintf("svr_export: %v\n\n\n", svr_export), 0)
			}
		}(j)
	}

	time.Sleep(5 * time.Second)

	ba.logger.Debug(fmt.Sprintf("Started all the replicas\n"), 0)

	clientOutputs := make([]string, num_clients)
	m := 0
	for j := int(num_replicas); j < int(num_replicas+num_clients); j++ {
		go func(i int, k int) {
			export_command := fmt.Sprintf("export LogFilePath=%vbench/logs/ RC_Ctrl=%v RC_Folder=%v/bench/ RC_LLevel=\"warn\" Rabia_ClosedLoop=false Rabia_NServers=%v Rabia_NFaulty=0 Rabia_NClients=%v Rabia_NConcurrency=1 Rabia_ClientBatchSize=%v Rabia_ClientTimeout=%v Rabia_ClientThinkTime=0 Rabia_ClientNRequests=0 Rabia_ClientArrivalRate=%v Rabia_ProxyBatchSize=%v Rabia_ProxyBatchTimeout=%v Rabia_NetworkBatchSize=0 Rabia_NetworkBatchTimeout=0 RC_Peers=%v Rabia_StorageMode=0",
				nodes[i].HomeDir, Controller, nodes[i].HomeDir, NServers, NClients, rabia_client_batch_size, duration, arrival_rate, rabia_proxy_batch_size, rabia_proxy_batch_timeout, RC_Peers_N)
			cli_export := fmt.Sprintf("export RC_Role=cli RC_Index=%v RC_Proxy=\"%v:11000\"", k, nodes[int64(i)-num_replicas].Ip)
			clientOutputs[i-int(num_replicas)] = nodes[i].ExecCmd(cli_export + ";" + export_command + ";" + "." + rabia_path)
			if i == int(num_replicas) {
				ba.logger.Debug(fmt.Sprintf("export_command: %v\n\n\n", export_command), 0)
				ba.logger.Debug(fmt.Sprintf("cli_export: %v\n\n\n", cli_export), 0)
			}
		}(j, m)
		m++
	}

	ba.logger.Debug(fmt.Sprintf("Started all the clients\n"), 0)

	crl_export := fmt.Sprintf("export RC_Role=ctrl")
	export_command := fmt.Sprintf("export LogFilePath=%vbench/logs/ RC_Ctrl=%v RC_Folder=%v/bench/ RC_LLevel=\"warn\" Rabia_ClosedLoop=false Rabia_NServers=%v Rabia_NFaulty=0 Rabia_NClients=%v Rabia_NConcurrency=1 Rabia_ClientBatchSize=%v Rabia_ClientTimeout=%v Rabia_ClientThinkTime=0 Rabia_ClientNRequests=0 Rabia_ClientArrivalRate=%v Rabia_ProxyBatchSize=%v Rabia_ProxyBatchTimeout=%v Rabia_NetworkBatchSize=0 Rabia_NetworkBatchTimeout=0 RC_Peers=%v Rabia_StorageMode=0",
		nodes[0].HomeDir, Controller, nodes[0].HomeDir, NServers, NClients, rabia_client_batch_size, duration, arrival_rate, rabia_proxy_batch_size, rabia_proxy_batch_timeout, RC_Peers_N)
	go nodes[0].ExecCmd(crl_export + ";" + export_command + ";" + "." + rabia_path)
	ba.logger.Debug(fmt.Sprintf("export_command: %v\n\n\n", export_command), 0)
	ba.logger.Debug(fmt.Sprintf("crl_export: %v\n\n\n", crl_export), 0)

	time.Sleep(5 * time.Second)

	ba.logger.Debug(fmt.Sprintf("Bootstrap complete\n"), 0)
	bootstrap_complete <- true

	time.Sleep(time.Duration(2*duration) * time.Second)

	ba.logger.Debug(fmt.Sprintf("Finished the clients\n"), 0)

	var wg1 sync.WaitGroup
	wg1.Add(int(num_replicas + num_clients))
	for j := 0; j < int(num_replicas+num_clients); j++ {
		go func(i int) {
			nodes[i].ExecCmd("pkill -KILL -f rabia")
			nodes[i].ExecCmd("pkill -KILL -f rabia")
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
	m = 0
	for j := int(num_replicas); j < int(num_replicas+num_clients); j++ {
		go func(i int, k int) {
			nodes[i].Get_Load(fmt.Sprintf("%vbench/logs/%v.txt", nodes[i].HomeDir, k), "logs/")
			wg2.Done()
		}(j, m)
		m++
	}
	wg2.Wait()
	ba.logger.Debug(fmt.Sprintf("Downloaded all the rabia client logs"), 0)

	command := "protocols/rabia/assets/performance_graph.py"
	file_names := []string{}

	m = 0
	for j := int(num_replicas); j < int(num_replicas+num_clients); j++ {
		logFile := filepath.Join(homeDir, fmt.Sprintf("toture-testing-consensus/logs/%v.txt", m))

		file_names = append(file_names, logFile)

		sshCmd = exec.Command("python3", []string{command, "rabia-" + strconv.Itoa(m), logFile}...)
		output, err = sshCmd.CombinedOutput()
		if err != nil {
			ba.logger.Debug(fmt.Sprintf("Error while generating performance graph "+err.Error()+" "+string(output)+"\n"), 0)
		} else {
			ba.logger.Debug(fmt.Sprintf("Rabia Performance graph generated successfully\n"+string(output)+"\n"), 0)
		}

		m++
	}

	sshCmd = exec.Command("python3", append([]string{command, "rabia"}, file_names...)...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug(fmt.Sprintf("Error while generating performance graphs "+err.Error()+" "+string(output)+"\n"), 0)
	} else {
		ba.logger.Debug(fmt.Sprintf("Rabia Performance graphs generated successfully\n"+string(output)+"\n"), 0)
	}

	result <- ba.GetPerformance(clientOutputs)
}

func (ba *Rabia) ExtractOptions(path string) protocols.ConsensusOptions {
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

	ba.logger.Debug(fmt.Sprintf("rabia options:\n %v\n", options.Option), 0)

	ba.options = options
	return options
}

func (ba *Rabia) GetPerformance(outputs []string) util.Performance {
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
