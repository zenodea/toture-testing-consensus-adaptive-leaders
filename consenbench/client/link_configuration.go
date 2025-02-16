package client

import (
	"fmt"
	"strconv"
	"strings"
	"toture-test/util"
)

type NetEmAttacker struct {
	Id                 int
	IP                 string
	Classid            string
	Handle             string
	NextNetEmCommands  [][]string
	DelayPackets       int
	LossPackets        int
	DuplicatePackets   int
	ReorderPackets     int
	CorruptPackets     int
	Rate               int
	logger             *util.Logger
	Ports_under_attack []string
	Device             string
	Prio               string
}

func (c *Client) NetInit(id_ip []string, ports_under_attack []string, device string) {
	RunCommand("tc", []string{"filter", "del", "dev", c.Options.Device}, c.logger)
	RunCommand("tc", []string{"qdisc", "del", "dev", c.Options.Device, "root"}, c.logger)
	RunCommand("tc", []string{"qdisc", "add", "dev", c.Options.Device, "root", "handle", "1:", "htb", "default", strconv.Itoa(len(id_ip) + 5)}, c.logger)
	RunCommand("tc", []string{"class", "add", "dev", c.Options.Device, "parent 1:", "classid 1:1", "htb", "rate", "1gbit", "ceil", "1gbit"}, c.logger)
	c.InitializeNetEmClients(id_ip, c.logger, ports_under_attack, device)
}

func (c *Client) InitializeNetEmClients(id_ip []string, logger *util.Logger, Ports_under_attack []string, device string) {
	c.Attacker.NetEmAttackers = make(map[int]*NetEmAttacker)
	for i := 0; i < len(id_ip); i++ {
		id, ip := strings.Split(id_ip[i], ":")[0], strings.Split(id_ip[i], ":")[1]
		id_int, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		c.Attacker.NetEmAttackers[id_int] = &NetEmAttacker{
			Id:                 id_int,
			IP:                 ip,
			Classid:            "1:" + strconv.Itoa((id_int+1)*10),
			Handle:             strconv.Itoa((id_int + 1) * 10),
			NextNetEmCommands:  [][]string{},
			DelayPackets:       0,
			LossPackets:        0,
			DuplicatePackets:   0,
			ReorderPackets:     0,
			CorruptPackets:     0,
			Rate:               10000000,
			logger:             logger,
			Ports_under_attack: Ports_under_attack,
			Device:             device,
			Prio:               strconv.Itoa((id_int + 1)),
		}

		RunCommand("tc", []string{"class", "add", "dev", c.Options.Device, "parent 1:1", "classid", c.Attacker.NetEmAttackers[id_int].Classid, "htb", "rate", "1gbit", "ceil", "1gbit", "prio", c.Attacker.NetEmAttackers[id_int].Prio}, c.logger)
		c.Attacker.NetEmAttackers[id_int].applyHandleToEachPort()
		debug := fmt.Sprintf("Initialized net em attacker with %v ", c.Attacker.NetEmAttackers[id_int])
		c.logger.Debug(debug, 5)
	}
}

// Execute the Pending commands

func (c *NetEmAttacker) ExecuteLastNetEmCommands() error {
	c.logger.Debug(fmt.Sprintf("Executing last netem commands %v", c.NextNetEmCommands), 3)
	for i := 0; i < len(c.NextNetEmCommands); i++ {
		if len(c.NextNetEmCommands[i]) == 0 {
			continue
		}
		RunCommand(c.NextNetEmCommands[i][0], c.NextNetEmCommands[i][1:], c.logger)
	}
	c.NextNetEmCommands = [][]string{}
	return nil
}

func (c *NetEmAttacker) SetNewHandler() error {
	if c.ReorderPackets > 0 && c.DelayPackets == 0 {
		c.DelayPackets = 1
		defer c.decrementDelay()
	}

	err := RunCommand("tc", []string{"qdisc", "add", "dev", c.Device, "parent", c.Classid, "handle", c.Handle + ":", "netem", "delay", strconv.Itoa(c.DelayPackets) + "ms", "loss", strconv.Itoa(c.LossPackets) + "%", "duplicate", strconv.Itoa(c.DuplicatePackets) + "%", "reorder", strconv.Itoa(c.ReorderPackets) + "%", "50%", "corrupt", strconv.Itoa(c.CorruptPackets) + "%", "rate", strconv.Itoa(c.Rate) + "kbit"}, c.logger)
	c.NextNetEmCommands = append(c.NextNetEmCommands, []string{"tc", "qdisc", "del", "dev", c.Device, "parent", c.Classid, "handle", c.Handle + ":", "netem", "delay", strconv.Itoa(c.DelayPackets) + "ms", "loss", strconv.Itoa(c.LossPackets) + "%", "duplicate", strconv.Itoa(c.DuplicatePackets) + "%", "reorder", strconv.Itoa(c.ReorderPackets) + "%", "50%", "corrupt", strconv.Itoa(c.CorruptPackets) + "%", "rate", strconv.Itoa(c.Rate) + "kbit"})
	c.logger.Debug("Set new net em handler", 3)
	return err
}

// apply the handle to each port

func (c *NetEmAttacker) applyHandleToEachPort() {

	for _, port := range c.Ports_under_attack {
		RunCommand("tc", []string{"filter", "add", "dev", c.Device, "protocol", "ip", "parent", "1:", "prio", c.Prio, "u32", "match", "ip", "dst", c.IP + "/32", "match", "ip", "dport", port, "0xffff", "flowid", c.Classid}, c.logger)
	}
}

func (c *NetEmAttacker) decrementDelay() {
	c.DelayPackets--
}

// set the delay

func (c *NetEmAttacker) SetDelay(f float32) {
	c.ExecuteLastNetEmCommands()
	c.DelayPackets = int(f)
	c.logger.Debug("set delay", 3)
	c.SetNewHandler()
}

// set the loss

func (c *NetEmAttacker) SetLoss(f float32) {
	c.ExecuteLastNetEmCommands()
	c.LossPackets = int(f)
	c.logger.Debug("set loss", 3)
	c.SetNewHandler()
}

// set the bandwidth

func (c *NetEmAttacker) SetBandwidth(f float32) {
	c.ExecuteLastNetEmCommands()
	c.Rate = int(f)
	c.logger.Debug("set bandwidth", 3)
	c.SetNewHandler()
}

func (c *Client) SetDelay(f float32, i int32) {
	node, ok := c.Attacker.NetEmAttackers[int(i)]
	if !ok {
		panic("NetEm Handler not found")
	}
	node.SetDelay(f)
}

func (c *Client) SetLoss(f float32, i int32) {
	node, ok := c.Attacker.NetEmAttackers[int(i)]
	if !ok {
		panic("NetEm handler not found" + strconv.Itoa(int(i)))
	}
	node.SetLoss(f)
}

func (c *Client) SetBandwidth(f float32, i int32) {
	node, ok := c.Attacker.NetEmAttackers[int(i)]
	if !ok {
		panic("Node not found")
	}
	node.SetBandwidth(f)
}

func (c *Client) CleanUp() error {
	RunCommand("tc", []string{"filter", "del", "dev", c.Attacker.Device}, c.logger)
	RunCommand("tc", []string{"qdisc", "del", "dev", c.Attacker.Device, "root"}, c.logger)
	return nil
}
