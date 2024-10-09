package controller

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"toture-test/protocols"
	baxos "toture-test/protocols/baxos"
	dedis_paxos "toture-test/protocols/dedis_paxos"
	dedis_raft "toture-test/protocols/dedis_raft"
	ping "toture-test/protocols/ping"
	rabia "toture-test/protocols/rabia"
	racs "toture-test/protocols/racs"
	sadl_racs "toture-test/protocols/sadl_racs"
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

// as you add more protocols, you need to add the protocol here

func (c *Controller) GetProtocolImpl(protocol string) protocols.Consensus {
	if protocol == "baxos" {
		return baxos.NewBaxos(c.logger)
	} else if protocol == "ping" {
		return ping.NewPing(c.logger)
	} else if protocol == "dedis_paxos" {
		return dedis_paxos.NewDedis_Paxos(c.logger)
	} else if protocol == "dedis_raft" {
		return dedis_raft.NewDedis_Raft(c.logger)
	} else if protocol == "rabia" {
		return rabia.NewRabia(c.logger)
	} else if protocol == "sadl_racs" {
		return sadl_racs.NewSadl_Racs(c.logger)
	} else if protocol == "racs" {
		return racs.NewRacs(c.logger)
	} else {
		panic("Unknown protocol")
	}
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

func (c *Controller) GetAttackImpl() Attack {
	switch c.Options.Attack {
	case "basic":
		return NewBasicAttack(c.logger)
	case "noop":
		return NewNoopAttack(c.logger)
	case "partition_1":
		return NewPartitionAttack_1(c.logger)
	case "partition_2":
		return NewPartitionAttack_2(c.logger)
	case "partition_3":
		return NewPartitionAttack_3(c.logger)
	case "partition_4":
		return NewPartitionAttack_4(c.logger)
	case "partition_5":
		return NewPartitionAttack_5(c.logger)
	case "partition_6":
		return NewPartitionAttack_6(c.logger)
	case "partition_7":
		return NewPartitionAttack_7(c.logger)
	case "delay_1":
		return NewDelayAttack_1(c.logger)
	case "delay_2":
		return NewDelayAttack_2(c.logger)
	case "delay_3":
		return NewDelayAttack_3(c.logger)
	case "delay_4":
		return NewDelayAttack_4(c.logger)
	case "bandwidth_1":
		return NewBandwidthAttack_1(c.logger)
	case "bandwidth_2":
		return NewBandwidthAttack_2(c.logger)
	case "bandwidth_3":
		return NewBandwidthAttack_3(c.logger)
	case "bandwidth_4":
		return NewBandwidthAttack_4(c.logger)
	case "skew_1":
		return NewSkewAttack_1(c.logger)
	case "skew_2":
		return NewSkewAttack_2(c.logger)
	case "skew_3":
		return NewSkewAttack_3(c.logger)
	case "straggler_1":
		return NewStragglerAttack_1(c.logger)
	case "straggler_2":
		return NewStragglerAttack_2(c.logger)
	case "straggler_3":
		return NewStragglerAttack_3(c.logger)
	case "crash_1":
		return NewCrashAttack_1(c.logger)
	case "crash_2":
		return NewCrashAttack_2(c.logger)
	case "crash_3":
		return NewCrashAttack_3(c.logger)
	case "drift_1":
		return NewDriftAttack_1(c.logger)
	case "drift_2":
		return NewDriftAttack_2(c.logger)
	case "drift_3":
		return NewDriftAttack_3(c.logger)
	default:
		panic("Unknown attack: " + c.Options.Attack)
	}
}
