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

type ZooKeeper struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewZooKeeper(logger *util.Logger) *ZooKeeper {
	return &ZooKeeper{
		logger: logger,
	}
}

func (ba *ZooKeeper) CopyConsensus(nodes []*common.Node) error {
	println("Running ZooKeeper")

	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("error while parsing num_replicas")
	}

	num_replicas_int, _ := strconv.Atoi(num_replicas)

	var wg sync.WaitGroup
	wg.Add(num_replicas_int)

	for i := 0; i < num_replicas_int; i++ {
		go func(j int) {
			nodes[j].ExecCmd(fmt.Sprintf("pkill -f QuorumPeerMain"))
			wg.Done()
		}(i)
	}
	wg.Wait()
	return nil
}

func (ba *ZooKeeper) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
	println("Running ZooKeeper Paxos")

	num_clients, ok := ba.options.Option["num_clients"]
	if !ok {
		panic("error while parsing num_clients")
	}

	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("error while parsing num_replicas")
	}

	tickTime, ok := ba.options.Option["tickTime"]
	if !ok {
		panic("error while parsing tickTime")
	}

	initLimit, ok := ba.options.Option["initLimit"]
	if !ok {
		panic("error while parsing initLimit")
	}

	syncLimit, ok := ba.options.Option["syncLimit"]
	if !ok {
		panic("error while parsing syncLimit")
	}

	num_clients_int, _ := strconv.Atoi(num_clients)
	num_replicas_int, _ := strconv.Atoi(num_replicas)

	outputs := make([]string, 0)
	outputMutex := &sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(num_replicas_int)

	ZOO_CFG_CONTENT := fmt.Sprintf("initLimit=1\ntickTime=%v\ndataDir=%vapache-zookeeper-3.8.1-bin/data\nclientPort=2181\ninitLimit=%v\nsyncLimit=%v\n", tickTime, nodes[0].HomeDir, initLimit, syncLimit)
	for i := 0; i < num_replicas_int; i++ {
		ZOO_CFG_CONTENT += fmt.Sprintf("server.%d=%v:2888:3888\n", i+1, nodes[i].Ip)
	}
	ZOO_CFG_CONTENT = ZOO_CFG_CONTENT[:len(ZOO_CFG_CONTENT)-1]

	for i := 0; i < num_replicas_int; i++ {
		go func(j int) {
			nodes[j].ExecCmd(fmt.Sprintf("pkill -f QuorumPeerMain"))
			nodes[j].ExecCmd(fmt.Sprintf("sudo rm -r -f %vapache-zookeeper*", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("wget https://archive.apache.org/dist/zookeeper/zookeeper-3.8.1/apache-zookeeper-3.8.1-bin.tar.gz -P %v", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("tar -xzf %vapache-zookeeper-3.8.1-bin.tar.gz", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("rm %vapache-zookeeper-3.8.1-bin.tar.gz", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vapache-zookeeper-3.8.1-bin/data", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("echo -e \"%v\" > %vapache-zookeeper-3.8.1-bin/conf/zoo.cfg", ZOO_CFG_CONTENT, nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("echo \"%v\" > %vapache-zookeeper-3.8.1-bin/data/myid", j+1, nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("rm -r %vbench/logs/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vbench/logs/", nodes[j].HomeDir))

			nodes[j].Put_Load(fmt.Sprintf("protocols/zoo_keeper/assets/client.py"), fmt.Sprintf("%v", nodes[j].HomeDir))

			go nodes[j].ExecCmd(fmt.Sprintf("%vapache-zookeeper-3.8.1-bin/bin/zkServer.sh start", nodes[j].HomeDir))

			time.Sleep(5 * time.Second)

			go func() {
				output := nodes[j].ExecCmd(fmt.Sprintf("python3 %vclient.py %v %v %v %v ", nodes[j].HomeDir, duration, num_clients_int, j, nodes[j].HomeDir+"bench/logs/"))
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
	fmt.Printf("bootstrap complete for zookeeper\n")
	time.Sleep(time.Duration(3*duration) * time.Second)
	fmt.Printf("finished running zookeeper\n")
	var wg1 sync.WaitGroup
	wg1.Add(num_replicas_int)
	for j := 0; j < num_replicas_int; j++ {
		go func(i int) {
			nodes[i].ExecCmd("pkill -f QuorumPeerMain")
			nodes[i].ExecCmd("pkill -f client.py")
			wg1.Done()
		}(j)
	}
	wg1.Wait()
	println("ZooKeeper killed")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("error getting home directory:" + err.Error())
	}

	sshCmd := exec.Command("rm", []string{"-r", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err := sshCmd.CombinedOutput()
	if err != nil {
		print("error while deleting logs/ " + err.Error() + " " + string(output) + "\n")
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
	fmt.Println("Downloaded all the zoo keepeter client logs")

	sshCmd = exec.Command("python3", []string{"protocols/zoo_keeper/assets/summary.py", filepath.Join(homeDir, "toture-testing-consensus/logs/"), filepath.Join(homeDir, "toture-testing-consensus/logs/")}...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		print("error while running summary " + err.Error() + " " + string(output) + "\n")
	} else {
		print("summary ran successfully\n" + string(output) + "\n")
	}

	p := ba.GetPerformance(outputs)
	result <- p
}

func (ba *ZooKeeper) ExtractOptions(path string) protocols.ConsensusOptions {
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

	fmt.Printf("ZooKeeper options:\n %v\n", options.Option)

	ba.options = options
	return options
}

func (ba *ZooKeeper) GetPerformance(outputs []string) util.Performance {
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
