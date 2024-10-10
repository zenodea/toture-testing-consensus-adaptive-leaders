package controller

import (
	"fmt"
	"os"
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

	var wg sync.WaitGroup
	wg.Add(len(c.Nodes))
	// copy the client binary to all the nodes
	for j := 0; j < len(c.Nodes); j++ {
		go func(i int) {
			c.Nodes[i].ExecCmd(fmt.Sprintf("sudo apt update"))
			c.Nodes[i].ExecCmd(fmt.Sprintf("sudo apt install iproute2"))
			c.Nodes[i].ExecCmd(fmt.Sprintf("sudo setcap cap_net_admin,cap_net_raw+ep $(which tc)"))
			c.Nodes[i].ExecCmd(fmt.Sprintf("getcap $(which tc)"))
			c.Nodes[i].ExecCmd(fmt.Sprintf("pkill -KILL -f bench"))
			c.Nodes[i].ExecCmd(fmt.Sprintf("rm -r %vbench", c.Nodes[i].HomeDir))
			c.Nodes[i].ExecCmd(fmt.Sprintf("mkdir -p %vbench", c.Nodes[i].HomeDir))
			c.Nodes[i].Put_Load("consenbench/bin/bench", fmt.Sprintf("%vbench/", c.Nodes[i].HomeDir))
			c.Nodes[i].Put_Load("consenbench/assets/ip.yaml", fmt.Sprintf("%vbench/", c.Nodes[i].HomeDir))
			wg.Done()
		}(j)
	}
	wg.Wait()

	fmt.Println("Copied the client binary to all the nodes")

	// start the client binary
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].Start_Client(c.Options.Device)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Started the client binary on all the nodes")

	// initiate the tcp connections
	c.NetworkInit()
	fmt.Println("Initialized the network layer with all clients")

	c.HandleClientMessages()

	time.Sleep(10 * time.Second)
	// close the clients
	c.CloseClients()
	fmt.Println("Closed the clients")
	c.DownloadClientLogs()
	fmt.Println("Downloaded the logs from the clients")
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
	// start the client binary
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].Start_Client(c.Options.Device)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Started the client binary on all the nodes")

	// initiate the tcp connections
	c.NetworkInit()
	fmt.Println("Initialized the network layer with all clients")

	c.HandleClientMessages()
	time.Sleep(5 * time.Second)

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
	fmt.Print("Bootstrap complete, starting attack from controller\n")

	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].StartUpdateStats()
	}

	time.Sleep(5 * time.Second)

	attack_impl.Attack(attackNodes, attackLinks, leaderOracle, c.Options.AttackDuration)
	fmt.Print("Attack complete\n")

	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].StopUpdateStats()
	}

	performance := <-performance_output_chan
	for key, value := range performance.Option {
		fmt.Printf("%v: %v\n", key, value)
	}

	c.PrintStats(int(num_replicas))

	c.CloseClients()
	fmt.Println("Closed the clients")
	c.DownloadClientLogs()
	fmt.Println("Downloaded the logs from the clients")
	fmt.Println("test complete")
}
