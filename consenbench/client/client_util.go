package client

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
	"toture-test/consenbench/common"
	"toture-test/util"
)

// Set the replica name and the ports that the replica is listening on

func (c *Client) InitAttacker(msg *common.ControlMsg) {
	replica_name := msg.StringArgs[0]
	ports := msg.StringArgs[1:]
	if len(ports) == 0 {
		panic("No ports provided")
	}
	c.Attacker.Process_name = replica_name
	c.Attacker.Ports_under_attack = ports
	c.logger.Debug("Set replica name to "+replica_name+" and ports to "+fmt.Sprintf("%v", ports), 3)
	c.Init(msg.Ips, ports, c.Attacker.Device)
}

// periodically send machine stats to the controller

func (c *Client) SendStats() {
	// send machine stats to the controller
	go func() {
		for true {
			cpu := util.GetCPUUsage()
			mem := util.GetMemoryUsage()
			packetsInRate, packetsOutRate := util.GetNetworkStats() // has a 1s sync delay

			perf_stats := []float32{float32(cpu), float32(mem), float32(packetsInRate), float32(packetsOutRate)}
			c.Network.Send(&common.RPCPairPeer{
				RpcPair: &common.RPCPair{
					Code: common.GetRPCCodes().ControlMsg,
					Obj: &common.ControlMsg{
						OperationType: int32(common.GetOperationCodes().Stats),
						FloatArgs:     perf_stats,
					},
				},
				Peer: c.ControllerId,
			})
			time.Sleep(100 * time.Millisecond)
			c.logger.Debug("Sent stats to controller", 0)

		}
	}()
}

// runCommand runs the given command with the provided arguments

func RunCommand(name string, arg []string, logger *util.Logger) error {
	// Create the command, running in a bash shell
	cmd := exec.Command("bash", "-c", name+" "+strings.Join(arg, " "))

	// Create pipes for stdout and stderr
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	// WaitGroup to track completion of goroutines
	var wg sync.WaitGroup

	// Start the command
	if err := cmd.Start(); err != nil {
		logger.Debug("Error running command:"+fmt.Sprintf("%v %v", name, arg)+" "+err.Error(), 5)
		return err
	}

	// Add two routines (stdout and stderr)
	wg.Add(2)

	// Goroutine to handle stdout
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, stdoutIn)
	}()

	// Goroutine to handle stderr
	go func() {
		defer wg.Done()
		io.Copy(os.Stderr, stderrIn)
	}()

	// Wait for the command to complete
	err := cmd.Wait()
	if err != nil {
		logger.Debug("Error running command:"+fmt.Sprintf("%v %v", name, arg)+" "+err.Error(), 5)
		return err
	}

	// Wait for both stdout and stderr goroutines to finish
	wg.Wait()

	// Log success message
	logger.Debug("Success command "+name+" "+strings.Join(arg, " "), 5)
	return nil
}

// Initialize the client

func (c *Client) Init(id_ip []string, ports_under_attack []string, device string) {
	c.NetInit(id_ip, ports_under_attack, device)
	c.logger.Debug("Initialized TC ", 3)
}
