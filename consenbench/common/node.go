package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/exec"
	"sync"
	"toture-test/util"
)

type NodeStat struct {
	cpu_usage   []float32
	mem_usage   []float32
	network_in  []float32
	network_out []float32
	update      bool
}

type Node struct {
	Id             int    `yaml:"Id"`
	Ip             string `yaml:"Ip"`
	Username       string `yaml:"Username"`
	HomeDir        string `yaml:"HomeDir"`
	PrivateKeyPath string `yaml:"privateKeyPath"`
	stat           *NodeStat
	statMutex      *sync.Mutex
	Logger         *util.Logger
}

func (n *Node) InitNode(logger *util.Logger) {
	n.stat = &NodeStat{
		cpu_usage:   make([]float32, 0),
		mem_usage:   make([]float32, 0),
		network_in:  make([]float32, 0),
		network_out: make([]float32, 0),
		update:      false,
	}
	n.statMutex = &sync.Mutex{}
	n.Logger = logger
}

// Execute a command on the node

func (n *Node) ExecCmd(cmd string) string {
	sshCmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-i", n.PrivateKeyPath, fmt.Sprintf("%s@%s", n.Username, n.Ip), cmd)
	output, err := sshCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("FAILED to execute %v via SSH, err:%v, output:%v for node:%v\n\n", fmt.Sprintf("%v", sshCmd), err, string(output), n.Id)
	} else {
		n.Logger.Debug(fmt.Sprintf("SUCCESS %v via SSH, output: %v for node: %v\n\n", fmt.Sprintf("%v", sshCmd), string(output), n.Id), 5)
	}
	return string(output)
}

// download the file from the remote location to the local location

func (n *Node) Get_Load(remote_location string, local_location string) error {

	scpCmd := exec.Command("scp", "-i", n.PrivateKeyPath, fmt.Sprintf("%s@%s:%s", n.Username, n.Ip, remote_location), local_location)
	output, err := scpCmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("FAILED to download file via SCP %v, error:%v, output:%v for node:%v\n\n", fmt.Sprintf("%v", scpCmd), err, string(output), n.Id))
	} else {
		n.Logger.Debug(fmt.Sprintf("SUCCESS download using %v for node %v\n\n", fmt.Sprintf("%v", scpCmd), n.Id), 0)
	}
	return nil
}

// upload the file from the local location to the remote location

func (n *Node) Put_Load(local_location string, remote_location string) error {
	scpCmd := exec.Command("scp", "-i", n.PrivateKeyPath, local_location, fmt.Sprintf("%s@%s:%s", n.Username, n.Ip, remote_location))
	output, err := scpCmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("FAILED to upload file via %v, err:%v, output:%s for node:%v\n\n", fmt.Sprintf("%v", scpCmd), err, string(output), n.Id))
	} else {
		n.Logger.Debug(fmt.Sprintf("SUCCESS Upload %v successful for node %v\n\n", fmt.Sprintf("%v", scpCmd), n.Id), 0)
	}
	return nil
}

// start client

func (n *Node) Start_Client(device string) error {

	fmt.Printf("Starting client on node: %v\n", n.Id)
	n.ExecCmd("pkill -KILL -f bench")
	n.ExecCmd("pkill -KILL -f fab")
	n.ExecCmd("tmux kill-server")
	exec.Command("pkill", "fab")

	if n.Logger.DebugOn {
		go n.ExecCmd(fmt.Sprintf("./bench/bench --node_config %vbench/ip.yaml --id %v --debug_on --debug_level %v --device %v ", n.HomeDir, n.Id, n.Logger.Level, device))
	} else {
		go n.ExecCmd(fmt.Sprintf("./bench/bench --node_config %vbench/ip.yaml --id %v --device %v", n.HomeDir, n.Id, device))
	}
	return nil
}

func (n *Node) StartUpdateStats() {
	n.statMutex.Lock()
	n.stat.update = true
	n.statMutex.Unlock()
}

func (n *Node) StopUpdateStats() {
	n.statMutex.Lock()
	n.stat.update = false
	n.statMutex.Unlock()
}

func (n *Node) UpdateStats(perf []float32) {
	n.statMutex.Lock()
	if n.stat.update {
		n.stat.cpu_usage = append(n.stat.cpu_usage, perf[0])
		n.stat.mem_usage = append(n.stat.mem_usage, perf[1])
		n.stat.network_in = append(n.stat.network_in, perf[2])
		n.stat.network_out = append(n.stat.network_out, perf[3])
	}
	n.statMutex.Unlock()
}

func (n *Node) GetStats() ([]float32, []float32, []float32, []float32) {
	n.statMutex.Lock()
	stats := NodeStat{
		cpu_usage:   GetNewArr(n.stat.cpu_usage),
		mem_usage:   GetNewArr(n.stat.mem_usage),
		network_in:  GetNewArr(n.stat.network_in),
		network_out: GetNewArr(n.stat.network_out),
	}
	n.statMutex.Unlock()
	return stats.cpu_usage, stats.mem_usage, stats.network_in, stats.network_out
}

func GetNewArr(usage []float32) []float32 {
	newArr := make([]float32, len(usage))
	for i := 0; i < len(usage); i++ {
		newArr[i] = usage[i]
	}
	return newArr
}

func GetNodes(filename string) []*Node {
	fmt.Printf("reading node data from filename: %v\n", filename)
	type Nodes struct {
		Nodes []*Node `yaml:"nodes"`
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("ERROR reading file: " + err.Error())
	}

	var nodes Nodes
	err = yaml.Unmarshal(data, &nodes)
	if err != nil {
		panic("ERROR unmarshalling YAML: " + err.Error())
	}

	client_nodes := nodes.Nodes[1:]
	return client_nodes

}

func GetController(filename string) *Node {
	type Nodes struct {
		Nodes []*Node `yaml:"nodes"`
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Error reading file: " + err.Error())
	}

	var nodes Nodes
	err = yaml.Unmarshal(data, &nodes)
	if err != nil {
		panic("Error unmarshalling YAML: " + err.Error())
	}

	controller := nodes.Nodes[0]
	return controller
}
