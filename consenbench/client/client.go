package client

import (
	"fmt"
	"sync"
	"time"
	"toture-test/consenbench/common"
	"toture-test/util"
)

type ClientOptions struct {
	NodeInfoFile string // the yaml file containing the ip address of each node, controller port, client port
	DebugOn      bool
	DebugLevel   int
	LogFilePath  string
	Device       string
}

type ClientAttacker struct {
	Logger             *util.Logger
	Ports_under_attack []string // ports under attack
	Process_name       string   // process under attack
	Device             string
	NetEmAttackers     map[int]*NetEmAttacker
}

type Client struct {
	Id           int
	Network      *common.Network // to communicate with the controller
	InputChan    chan *common.RPCPairPeer
	logger       *util.Logger
	Options      ClientOptions
	ControllerId int
	Attacker     *ClientAttacker
}

func NewClient(Id int, options ClientOptions) *Client {
	c := &Client{
		Id:        Id,
		InputChan: make(chan *common.RPCPairPeer, 10000),
		logger:    util.NewLogger(options.DebugLevel, options.DebugOn, options.LogFilePath),
		Options:   options,
	}

	attacker := &ClientAttacker{
		Logger: c.logger,
		Device: options.Device,
	}

	c.Attacker = attacker
	return c
}

// initialize the network layer and run the client

func (c *Client) Run() {

	nodes := common.GetNodes(c.Options.NodeInfoFile)
	listenAddress := ""
	for i := 0; i < len(nodes); i++ {
		if nodes[i].Id == c.Id {
			listenAddress = nodes[i].Ip + ":10080"
		}
	}
	if listenAddress == "" {
		panic("node id not found in the node info file")
	}

	controller := common.GetController(c.Options.NodeInfoFile)
	c.ControllerId = controller.Id

	network_config := common.NetworkConfig{
		ListenAddress:   listenAddress,
		RemoteAddresses: map[int]string{controller.Id: controller.Ip + ":10080"},
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
	c.logger.Debug("Connected to controller, both ways!", 0)

	c.SendStats()

	c.run()

	for true {
		time.Sleep(5 * time.Second)
	}
}

// respond to different messages from the controller

func (c *Client) run() {
	go func() {
		for true {
			rpcPair := <-c.InputChan
			c.logger.Debug(fmt.Sprintf("%v", rpcPair.RpcPair), 0)
			switch rpcPair.RpcPair.Code {
			case common.GetRPCCodes().ControlMsg:
				// handle control message
				ctrlMsg := rpcPair.RpcPair.Obj.(*common.ControlMsg)
				c.Handle(ctrlMsg)
			default:
				c.logger.Debug("Unknown message type", 0)
			}
		}
	}()
}
