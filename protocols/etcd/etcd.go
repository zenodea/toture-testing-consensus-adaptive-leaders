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

type ETCD struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewETCD(logger *util.Logger) *ETCD {
	return &ETCD{
		logger: logger,
	}
}

func (ba *ETCD) CopyConsensus(nodes []*common.Node) error {
	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("error while parsing num_replicas")
	}

	num_replicas_int, _ := strconv.Atoi(num_replicas)

	var wg sync.WaitGroup
	wg.Add(num_replicas_int)

	ETCD_VERSION := "v3.5.9"
	ETCD_DOWNLOAD_URL := fmt.Sprintf("https://github.com/etcd-io/etcd/releases/download/%v/etcd-%v-linux-amd64.tar.gz", ETCD_VERSION, ETCD_VERSION)

	for i := 0; i < num_replicas_int; i++ {
		go func(j int) {
			nodes[j].ExecCmd(fmt.Sprintf("sudo pkill -f etcd"))
			nodes[j].ExecCmd(fmt.Sprintf("sudo rm -f /usr/local/bin/etcd;sudo rm -f /usr/local/bin/etcdctl; sudo rm -f /usr/local/bin/etcdutl"))
			nodes[j].ExecCmd(fmt.Sprintf("sudo systemctl stop etcd; sudo systemctl disable etcd; sudo rm -f /etc/systemd/system/etcd.service"))
			nodes[j].ExecCmd(fmt.Sprintf("sudo rm -rf %vetcd/; sudo rm -r /var/lib/etcd; sudo rm -r /tmp/etcd", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("sudo rm -rf %vinfra%v.etcd/", nodes[j].HomeDir, j))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vetcd/data", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("wget -q %v && tar xzvf etcd-%v-linux-amd64.tar.gz", ETCD_DOWNLOAD_URL, ETCD_VERSION))
			nodes[j].ExecCmd(fmt.Sprintf("mv etcd-%v-linux-amd64/etcd* %vetcd/", ETCD_VERSION, nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("rm -rf etcd-%v-linux-amd64*", ETCD_VERSION))
			nodes[j].Put_Load("protocols/etcd/assets/client.py", fmt.Sprintf("%vetcd/client.py", nodes[j].HomeDir))
			wg.Done()
		}(i)
	}
	wg.Wait()
	return nil
}

func (ba *ETCD) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
	ba.logger.Debug(fmt.Sprintf("Running ETCD Raft"), 0)

	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("error while parsing num_replicas")
	}

	view_timeout, ok := ba.options.Option["view_timeout"]
	if !ok {
		panic("error while parsing view_timeout")
	}

	num_replicas_int, _ := strconv.Atoi(num_replicas)

	num_clients, ok := ba.options.Option["num_clients"]
	if !ok {
		panic("error while parsing num_clients")
	}

	num_clients_int, _ := strconv.Atoi(num_clients)

	outputs := make([]string, 0)
	outputMutex := &sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(num_replicas_int)

	CLUSTER := ""
	for i := 0; i < num_replicas_int; i++ {
		CLUSTER += fmt.Sprintf("infra%d=http://%v:2380,", i, nodes[i].Ip)
	}
	CLUSTER = CLUSTER[:len(CLUSTER)-1]

	for i := 0; i < num_replicas_int; i++ {
		go func(j int) {
			nodes[j].ExecCmd(fmt.Sprintf("rm -r %vbench/logs/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vbench/logs/", nodes[j].HomeDir))
			go nodes[j].ExecCmd(fmt.Sprintf("%vetcd/etcd --log-level error --name infra%v --initial-advertise-peer-urls http://%v:2380 --listen-peer-urls http://%v:2380 --listen-client-urls http://%v:2379,http://127.0.0.1:2379 --advertise-client-urls http://%v:2379 --initial-cluster-token etcd-cluster-1 --initial-cluster %v --initial-cluster-state new --election-timeout %v", nodes[j].HomeDir, j, nodes[j].Ip, nodes[j].Ip, nodes[j].Ip, nodes[j].Ip, CLUSTER, view_timeout))

			time.Sleep(5 * time.Second)

			go func() {
				output := nodes[j].ExecCmd(fmt.Sprintf("python3 %vetcd/client.py %v %v %v %v", nodes[j].HomeDir, duration, num_clients_int, j, nodes[j].HomeDir+"bench/logs/"))
				outputMutex.Lock()
				outputs = append(outputs, output)
				outputMutex.Unlock()
			}()

			time.Sleep(5 * time.Second)
			wg.Done()
		}(i)
	}
	wg.Wait()
	bootstrap_complete <- true
	ba.logger.Debug(fmt.Sprintf("bootstrap complete for etcd\n"), 0)
	time.Sleep(time.Duration(3*duration) * time.Second)
	ba.logger.Debug(fmt.Sprintf("finished running etcd"), 0)
	var wg1 sync.WaitGroup
	wg1.Add(num_replicas_int)
	for j := 0; j < num_replicas_int; j++ {
		go func(i int) {
			nodes[i].ExecCmd("pkill -f etcd")
			nodes[i].ExecCmd("pkill -f client.py")
			wg1.Done()
		}(j)
	}
	wg1.Wait()
	ba.logger.Debug(fmt.Sprintf("ETCD Raft killed"), 0)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("error getting home directory:" + err.Error())
	}

	sshCmd := exec.Command("rm", []string{"-r", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err := sshCmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug(fmt.Sprintf("error while deleting logs/ "+err.Error()+" "+string(output)+"\n"), 0)
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
	wg2.Add(num_replicas_int)

	for j := 0; j < num_replicas_int; j++ {
		go func(i int) {
			for k := 0; k < num_clients_int; k++ {
				nodes[i].Get_Load(fmt.Sprintf("%vbench/logs/%v_%v.log", nodes[i].HomeDir, i, k), "logs/")
			}
			wg2.Done()
		}(j)

	}
	wg2.Wait()
	ba.logger.Debug(fmt.Sprintf("Downloaded all the etcd client logs"), 0)

	sshCmd = exec.Command("python3", []string{"protocols/etcd/assets/summary.py", filepath.Join(homeDir, "toture-testing-consensus/logs/"), filepath.Join(homeDir, "toture-testing-consensus/logs/")}...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug(fmt.Sprintf("error while running summary "+err.Error()+" "+string(output)+"\n"), 0)
	} else {
		ba.logger.Debug(fmt.Sprintf("summary ran successfully\n"+string(output)+"\n"), 0)
	}

	p := ba.GetPerformance(outputs)
	result <- p
}

func (ba *ETCD) ExtractOptions(path string) protocols.ConsensusOptions {
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

	ba.logger.Debug(fmt.Sprintf("ETCD options:\n %v\n", options.Option), 0)

	ba.options = options
	return options
}

func (ba *ETCD) GetPerformance(outputs []string) util.Performance {
	sum_throughput := 0.0
	sum_latency := 0.0
	entries := 0.0
	for _, output := range outputs {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) > 0 && strings.Contains(line, "perf") && len(strings.Split(line, ",")) >= 3 {
				lat, tx := strings.Split(line, ",")[1], strings.Split(line, ",")[2]
				latency, _ := strconv.ParseFloat(lat, 64)
				throughput, _ := strconv.ParseFloat(tx, 64)
				sum_latency += latency
				sum_throughput += throughput
				entries = entries + 1
				break
			}
		}
	}

	p := util.Performance{
		map[string]string{"latency": fmt.Sprintf("%v", sum_latency/entries), "throughput": fmt.Sprintf("%v", sum_throughput)}}

	fmt.Printf("%v,%v,%v,", sum_throughput, sum_latency/entries, 0)

	return p
}
