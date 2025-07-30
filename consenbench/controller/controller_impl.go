package controller

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"toture-test/protocols"
	"toture-test/util"
)

// setup machines

func (c *Controller) SetupMachines() error {

	c.InitiliazeNodes()

	c.logger.Debug(fmt.Sprintf("Setting up machines"), 0)

	var wg sync.WaitGroup
	wg.Add(len(c.Nodes))
	for j := 0; j < len(c.Nodes); j++ {
		go func(i int) {

			cmd := "set -e;sudo apt -y update;sudo apt-mark hold grub-efi-amd64 grub-pc grub-common; sudo apt -y upgrade;sudo apt-get -y autoremove; sudo apt install -y git iproute2 cmake build-essential clang pkg-config libssl-dev python3 python3-pip openjdk-11-jdk;sudo rm -rf /usr/local/go; wget -q https://go.dev/dl/go1.19.linux-amd64.tar.gz;sudo tar -C /usr/local -xzf go1.19.linux-amd64.tar.gz;rm go1.19.linux-amd64.tar.gz;echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc;. ~/.bashrc;export PATH=$PATH:/usr/local/go/bin:$HOME/.cargo/bin; sudo rm -rf ~/.cargo ~/.rustup;curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y;. $HOME/.cargo/env;rustup default stable; sudo chmod u+s /usr/sbin/tc;sudo setcap cap_net_admin,cap_net_raw+ep $(which tc);getcap $(which tc);pip3 install --user --no-cache-dir --force-reinstall fabric invoke decorator boto3 etcd3 kazoo protobuf==3.19.6; sudo apt install -y python3-matplotlib"
			c.Nodes[i].ExecCmd(fmt.Sprintf(cmd))
			wg.Done()
		}(j)
	}
	wg.Wait()

	c.logger.Debug(fmt.Sprintln("done setup all nodes"), 0)
	return nil
}

// copy clients and start the client binary, check connections and close clients

func (c *Controller) BootstrapClients() error {

	c.InitiliazeNodes()

	c.logger.Debug(fmt.Sprintf("Copying the client binary to all nodes"), 0)

	process_names := []string{"node", "mysticeti", "replica", "replica", "epaxos_server", "etcd", "node", "node", "mysticeti", "mysticeti", "dummy", "replica", "rabia", "replica", "replica", "node", "QuorumPeerMain"}

	var wg sync.WaitGroup
	wg.Add(len(c.Nodes))
	// copy the client binary to all the nodes
	for j := 0; j < len(c.Nodes); j++ {
		go func(i int) {

			for p := 0; p < len(process_names); p++ {
				c.Nodes[i].ExecCmd(fmt.Sprintf("pkill -KILL -f " + process_names[p]))
			}
			c.Nodes[i].ExecCmd(fmt.Sprintf("pkill -KILL -f bench"))
			c.Nodes[i].ExecCmd(fmt.Sprintf("pkill -KILL -f fab"))
			c.Nodes[i].ExecCmd(fmt.Sprintf("rm -r %vbench", c.Nodes[i].HomeDir))
			c.Nodes[i].ExecCmd(fmt.Sprintf("mkdir -p %vbench", c.Nodes[i].HomeDir))
			c.Nodes[i].Put_Load("consenbench/bin/bench", fmt.Sprintf("%vbench/", c.Nodes[i].HomeDir))
			c.Nodes[i].Put_Load("consenbench/assets/ip.yaml", fmt.Sprintf("%vbench/", c.Nodes[i].HomeDir))
			c.logger.Debug(fmt.Sprintf("done copying the client binary to %v", i), 0)
			wg.Done()
		}(j)
	}
	wg.Wait()

	c.logger.Debug(fmt.Sprintln("Copied the client binary to all the nodes"), 0)

	// start the client binary
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].Start_Client(c.Options.Device)
	}
	time.Sleep(5 * time.Second)
	c.logger.Debug(fmt.Sprintf("Started the client binary on all the nodes"), 0)

	// initiate the tcp connections
	c.NetworkInit()
	c.logger.Debug(fmt.Sprintf("Initialized the network layer with all clients"), 0)

	c.HandleClientMessages()

	time.Sleep(10 * time.Second)
	// close the clients
	c.CloseClients()
	c.logger.Debug(fmt.Sprintf("Closed the clients"), 0)
	c.DownloadClientLogs()
	c.logger.Debug(fmt.Sprintf("Downloaded the logs from the clients"), 0)
	os.Exit(0)
	return nil

}

// copy the consensus binary

func (c *Controller) CopyConsensus(protocol string) {
	c.InitiliazeNodes()
	var protocol_impl protocols.Consensus
	protocol_impl = c.GetProtocolImpl(protocol)
	protocol_impl.ExtractOptions("protocols/" + protocol + "/assets/options.yaml")
	protocol_impl.CopyConsensus(c.Nodes)
}

// run the controller

func (c *Controller) Run(protocol string) {
	c.InitiliazeNodes()

	var wg2 sync.WaitGroup
	wg2.Add(len(c.Nodes))

	for i := 0; i < len(c.Nodes); i++ {
		go func(j int) {
			c.Nodes[j].ExecCmd(fmt.Sprintf("rm -rf %va_mysticeti", c.Nodes[j].HomeDir))
			c.Nodes[j].ExecCmd(fmt.Sprintf("rm -rf %vasync-mystecity", c.Nodes[j].HomeDir))
			c.Nodes[j].ExecCmd(fmt.Sprintf("rm -rf %vresults", c.Nodes[j].HomeDir))
			c.Nodes[j].ExecCmd(fmt.Sprintf("find %v -type f -name \"stable-store*\" -delete", c.Nodes[j].HomeDir))
			c.Nodes[j].ExecCmd(fmt.Sprintf("find %v -type f -name \"client-times*\" -delete", c.Nodes[j].HomeDir))
			c.Nodes[j].ExecCmd(fmt.Sprintf("find %v -type f -name \"db-*\" -delete", c.Nodes[j].HomeDir))
			c.Nodes[j].ExecCmd(fmt.Sprintf("find %v -type f -name \"storage-*\" -delete", c.Nodes[j].HomeDir))
			c.Nodes[j].ExecCmd(fmt.Sprintf("df -h"))
			c.logger.Debug(fmt.Sprintf("done cleaning up the node %v", j), 0)
			wg2.Done()
		}(i)
	}

	wg2.Wait()

	fmt.Printf("%v,%v,", protocol, c.Options.Attack)

	// start the client binary
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].Start_Client(c.Options.Device)
	}
	time.Sleep(5 * time.Second)
	c.logger.Debug(fmt.Sprintf("Started the client binary on all the nodes"), 0)

	// initiate the tcp connections
	c.NetworkInit()
	c.logger.Debug(fmt.Sprintf("Initialized the network layer with all clients"), 0)

	c.HandleClientMessages()
	time.Sleep(5 * time.Second)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("error getting home directory:" + err.Error())
	}

	cmd := exec.Command("rm", []string{"-r", filepath.Join(homeDir, "toture-testing-consensus/final-results/"+protocol+"/"+c.Options.Attack)}...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.logger.Debug(fmt.Sprintf("Error while deleting final-results/"+protocol+"/"+c.Options.Attack+err.Error()+" "+string(output)+"\n"), 0)
	} else {
		c.logger.Debug(fmt.Sprintf("deleted final-results sub directory successfully\n"+string(output)+"\n"), 0)
	}

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

	cmd = exec.Command("mkdir", []string{"-p", filepath.Join(homeDir, "toture-testing-consensus/final-results/"+param_load+"/"+param_size+"/"+protocol+"/"+c.Options.Attack)}...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		panic("Error while creating final-results/" + protocol + "/" + c.Options.Attack + err.Error() + " " + string(output) + "\n")
	} else {
		c.logger.Debug(fmt.Sprintf("created final-results/ sub directory successfully\n"+string(output)+"\n"), 0)
	}

	cmd = exec.Command("rm", []string{"-r", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		c.logger.Debug(fmt.Sprintf("error while deleting logs/ "+err.Error()+" "+string(output)+"\n"), 0)
	} else {
		c.logger.Debug(fmt.Sprintf("deleted local logs/ successfully\n"+string(output)+"\n"), 0)
	}

	cmd = exec.Command("mkdir", []string{"-p", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		panic("Error while creating logs/ " + err.Error() + " " + string(output) + "\n")
	} else {
		c.logger.Debug(fmt.Sprintf("created logs/ successfully\n"+string(output)+"\n"), 0)
	}

	protocol_impl := c.GetProtocolImpl(protocol)
	options := protocol_impl.ExtractOptions("protocols/" + protocol + "/assets/options.yaml")

	bootstrap_complete_chan := make(chan bool)
	performance_output_chan := make(chan util.Performance)

	go protocol_impl.Bootstrap(c.Nodes, c.Options.AttackDuration, performance_output_chan, bootstrap_complete_chan)

	num_replicas, err := strconv.ParseInt(options.Option["num_replicas"], 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}
	process_name, ok := options.Option["process_name"]
	if !ok {
		panic("error while parsing process name")

	}

	ports, ok := options.Option["ports"]
	if !ok {
		panic(err.Error() + " while parsing ports")
	}

	attackNodes, attackLinks, leaderOracle := GetAttackObjects(int(num_replicas), process_name, c.Nodes, c, c.logger, strings.Split(ports, ","))
	attack_impl := c.GetAttackImpl()

	<-bootstrap_complete_chan // wait for the bootstrap to complete
	c.logger.Debug(fmt.Sprintf("Bootstrap complete, starting attack from controller\n"), 0)

	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].StartUpdateStats()
	}

	attack_impl.Attack(attackNodes, attackLinks, leaderOracle, c.Options.AttackDuration)
	c.logger.Debug(fmt.Sprint("Attack complete\n"), 0)

	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].StopUpdateStats()
	}

	perf := <-performance_output_chan
	c.logger.Debug(fmt.Sprintf("Performance %v \n", perf), 0)

	c.PrintStats(int(num_replicas))

	c.CloseClients()
	c.logger.Debug(fmt.Sprintf("Closed the clients"), 0)
	c.DownloadClientLogs()
	c.logger.Debug(fmt.Sprintf("Downloaded the logs from the clients"), 0)

	logDir := filepath.Join(homeDir, "toture-testing-consensus", "logs")
	destDir := filepath.Join(homeDir, "toture-testing-consensus", "final-results", param_load, param_size, protocol, c.Options.Attack)

	cmd = exec.Command("sh", "-c", "mv "+logDir+"/*.pdf "+destDir)

	output, err = cmd.CombinedOutput()
	if err != nil {
		panic("Error while moving to final-results/" + protocol + "/" + c.Options.Attack + err.Error() + " " + string(output) + "\n")
	} else {
		c.logger.Debug(fmt.Sprintf("moved final-results/ sub directory successfully\n"+string(output)+"\n"), 0)
	}

	logDir = filepath.Join(homeDir, "toture-testing-consensus", "bench", "log.log")

	cmd = exec.Command("sh", "-c", "mv "+logDir+"  "+destDir)

	output, err = cmd.CombinedOutput()
	if err != nil {
		panic("Error while moving to final-results/" + protocol + "/" + c.Options.Attack + err.Error() + " " + string(output) + "\n")
	} else {
		c.logger.Debug(fmt.Sprintf("moved final-results/ sub directory successfully\n"+string(output)+"\n"), 0)
	}

	var wg sync.WaitGroup
	wg.Add(int(num_replicas))

	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			c.Nodes[j].ExecCmd(fmt.Sprintf("rm -r %vbench/logs/", c.Nodes[j].HomeDir))
			c.Nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vbench/logs/", c.Nodes[j].HomeDir))
			c.Nodes[j].ExecCmd(fmt.Sprintf("df -h"))
			wg.Done()
		}(i)
	}

	wg.Wait()
	c.logger.Debug(fmt.Sprintf("test complete"), 0)
}
