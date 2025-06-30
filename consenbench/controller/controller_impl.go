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

// copy clients and start the client binary, check connections and close clients

func (c *Controller) BootstrapClients() error {

	c.InitiliazeNodes()

	c.logger.Debug(fmt.Sprintf("Copying the client binary to all nodes"), 0)

	var wg sync.WaitGroup
	wg.Add(len(c.Nodes))
	// copy the client binary to all the nodes
	for j := 0; j < len(c.Nodes); j++ {
		go func(i int) {
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
	fmt.Printf("%v,%v,", protocol, c.Options.Attack)
	c.InitiliazeNodes()
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

	cmd = exec.Command("mkdir", []string{"-p", filepath.Join(homeDir, "toture-testing-consensus/final-results/"+protocol+"/"+c.Options.Attack)}...)
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

	_ = <-performance_output_chan

	c.PrintStats(int(num_replicas))

	c.CloseClients()
	c.logger.Debug(fmt.Sprintf("Closed the clients"), 0)
	c.DownloadClientLogs()
	c.logger.Debug(fmt.Sprintf("Downloaded the logs from the clients"), 0)

	logDir := filepath.Join(homeDir, "toture-testing-consensus", "logs")
	destDir := filepath.Join(homeDir, "toture-testing-consensus", "final-results", protocol, c.Options.Attack)

	cmd = exec.Command("sh", "-c", "mv "+logDir+"/*.pdf "+destDir)

	output, err = cmd.CombinedOutput()
	if err != nil {
		panic("Error while moving to final-results/" + protocol + "/" + c.Options.Attack + err.Error() + " " + string(output) + "\n")
	} else {
		c.logger.Debug(fmt.Sprintf("moved final-results/ sub directory successfully\n"+string(output)+"\n"), 0)
	}

	c.logger.Debug(fmt.Sprintf("test complete"), 0)
}
