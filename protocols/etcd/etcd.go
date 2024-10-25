package consensus

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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
	return nil
}

func (ba *ETCD) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
	println("Running ETCD Raft")

	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("error while parsing num_replicas")
	}

	num_replicas_int, _ := strconv.Atoi(num_replicas)

	outputs := make([]string, 0)
	outputMutex := &sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(num_replicas_int)

	ETCD_VERSION := "v3.5.9"
	ETCD_DOWNLOAD_URL := fmt.Sprintf("https://github.com/etcd-io/etcd/releases/download/%v/etcd-%v-linux-amd64.tar.gz", ETCD_VERSION, ETCD_VERSION)
	CLUSTER := ""
	for i := 0; i < num_replicas_int; i++ {
		CLUSTER += fmt.Sprintf("infra%d=http://%v:2380,", i, nodes[i].Ip)
	}
	CLUSTER = CLUSTER[:len(CLUSTER)-1]

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
			nodes[j].ExecCmd(fmt.Sprintf("sudo apt update ; sudo apt install -y python3 python3-pip"))
			nodes[j].ExecCmd(fmt.Sprintf("pip3 install etcd3"))
			nodes[j].ExecCmd(fmt.Sprintf("pip3 install protobuf==3.19.6"))
			go nodes[j].ExecCmd(fmt.Sprintf("%vetcd/etcd --log-level error --name infra%v --initial-advertise-peer-urls http://%v:2380 --listen-peer-urls http://%v:2380 --listen-client-urls http://%v:2379,http://127.0.0.1:2379 --advertise-client-urls http://%v:2379 --initial-cluster-token etcd-cluster-1 --initial-cluster %v --initial-cluster-state new", nodes[j].HomeDir, j, nodes[j].Ip, nodes[j].Ip, nodes[j].Ip, nodes[j].Ip, CLUSTER))
			nodes[j].Put_Load("protocols/etcd/assets/client.py", fmt.Sprintf("%vetcd/client.py", nodes[j].HomeDir))
			for k := 0; k < 7; k++ {
				go func() {
					output := nodes[j].ExecCmd(fmt.Sprintf("python3 %vetcd/client.py %v", nodes[j].HomeDir, duration))
					outputMutex.Lock()
					outputs = append(outputs, output)
					outputMutex.Unlock()
				}()
			}
			time.Sleep(5 * time.Second)
			wg.Done()
		}(i)
	}
	wg.Wait()
	bootstrap_complete <- true
	fmt.Printf("bootstrap complete for etcd\n")
	time.Sleep(time.Duration(2*duration) * time.Second)
	fmt.Printf("finished running etcd")
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
	println("ETCD Raft killed")
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

	fmt.Printf("ETCD options:\n %v\n", options.Option)

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

	return util.Performance{
		map[string]string{"latency": fmt.Sprintf("%v", sum_latency/entries), "throughput": fmt.Sprintf("%v", sum_throughput)}}
}
