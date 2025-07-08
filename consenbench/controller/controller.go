package controller

import (
	"fmt"
	"sync"
	"toture-test/consenbench/common"
	"toture-test/util"
)

type ControllerOptions struct {
	AttackDuration int    // second
	Attack         string // set of attacks to run
	NodeInfoFile   string // the yaml file containing the ip address of each node, controller port, client port
	DebugOn        bool
	DebugLevel     int
	LogFilePath    string
	Device         string
}

// Controller struct
type Controller struct {
	Id        int
	Nodes     []*common.Node
	Network   *common.Network          // to communicate with the clients
	InputChan chan *common.RPCPairPeer // message from the clients
	Options   ControllerOptions
	logger    *util.Logger
}

func NewController(Id int, Options ControllerOptions) *Controller {
	return &Controller{
		Id:        Id,
		InputChan: make(chan *common.RPCPairPeer, 10000),
		Options:   Options,
		logger:    util.NewLogger(Options.DebugLevel, Options.DebugOn, Options.LogFilePath),
	}
}

// start listening to incoming connections and connect to all remote nodes

func (c *Controller) NetworkInit() error {
	network_config := common.NetworkConfig{
		ListenAddress:   common.GetController(c.Options.NodeInfoFile).Ip + ":10080",
		RemoteAddresses: common.GetRemoteAddresses(c.Nodes),
	}
	c.Network = common.NewNetwork(c.Id, &network_config, c.InputChan, c.logger)
	c.Network.RegisterRPC(&common.ControlMsg{}, common.GetRPCCodes().ControlMsg)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		c.Network.ConnectRemotes()
		wg.Done()
	}()
	c.Network.Listen()
	wg.Wait()
	c.logger.Debug("Connected to all remote nodes, both ways!", 3)
	return nil
}

// handle the messages from the clients about machine stats

func (c *Controller) HandleClientMessages() error {
	go func() {
		for true {
			msg := <-c.InputChan
			c.logger.Debug(fmt.Sprintf("received from client %v %v ", msg.Peer, msg.RpcPair), -1)
			switch msg.RpcPair.Code {
			case common.GetRPCCodes().ControlMsg:
				// handle control message
				ctrlMsg := msg.RpcPair.Obj.(*common.ControlMsg)
				c.Handle(ctrlMsg, msg.Peer)
			default:
				c.logger.Debug("Unknown message type", 0)
			}
		}
	}()
	return nil
}

func (c *Controller) InitiliazeNodes() {
	c.Nodes = common.GetNodes(c.Options.NodeInfoFile)
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].InitNode(c.logger)
	}
}

func (c *Controller) CloseClients() {
	c.Network.Broadcast(&common.RPCPair{
		Code: common.GetRPCCodes().ControlMsg,
		Obj: &common.ControlMsg{
			OperationType: int32(common.GetOperationCodes().ShutDown),
		},
	})
}

func (c *Controller) DownloadClientLogs() {
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].Get_Load(c.Nodes[i].HomeDir+"bench/log.log", fmt.Sprintf("bench/%v.log", c.Nodes[i].Id))
	}
}

func (c *Controller) PrintStats(num_replicas int) {
	cpu_usage, mem_usage, network_in, network_out := float32(0.0), float32(0.0), float32(0.0), float32(0.0)
	counter := 0
	for i := 0; i < num_replicas; i++ {
		cpu, mem, net_in, net_out := c.Nodes[i].GetStats()
		c.logger.Debug(fmt.Sprintf("Node: %v -- CPU: %v, Memory: %v, Network In: %v, Network Out: %v", c.Nodes[i].Id, cpu, mem, net_in, net_out), 3)
		cpu_usage += Sum(cpu)
		mem_usage += Sum(mem)
		network_in += Sum(net_in)
		network_out += Sum(net_out)
		counter += len(cpu)
	}
	fmt.Printf("%v,%v,%v,%v\n", cpu_usage/float32(counter), mem_usage/float32(counter), network_in/float32(counter), network_out/float32(counter))
}
