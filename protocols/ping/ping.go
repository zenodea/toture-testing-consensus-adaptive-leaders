package consensus

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"sync"
	"time"
	"toture-test/consenbench/common"
	"toture-test/protocols"
	"toture-test/util"
)

type Ping struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewPing(logger *util.Logger) *Ping {
	return &Ping{
		logger: logger,
	}
}

func (ba *Ping) CopyConsensus(nodes []*common.Node) error {
	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("num_replicas not found in options")
	}

	num_replicas_int, err := strconv.ParseInt(num_replicas, 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}

	if num_replicas_int > int64(len(nodes)) {
		panic("Not enough nodes to deploy Ping")
	}

	config_inputs := []string{"protocols/ping/assets/config-generate.py", num_replicas}
	for i := int64(0); i < num_replicas_int; i++ {
		config_inputs = append(config_inputs, nodes[i].Ip)
	}

	ba.logger.Debug(fmt.Sprintf("Running python command: %v\n", config_inputs), 0)

	sshCmd := exec.Command("python3", config_inputs...)
	output, err := sshCmd.CombinedOutput()
	if err != nil {

		panic("Error while running config-generate.py " + err.Error() + " " + string(output))
	}
	if len(output) > 0 {
		// write the output to protocols/ping/assets/ip_config.yaml
		err = ioutil.WriteFile("protocols/ping/assets/ip_config.yaml", output, 0644)
		if err != nil {
			panic("Error while writing to ip_config.yaml " + err.Error())
		} else {
			ba.logger.Debug(fmt.Sprintf("ip_config.yaml written successfully with content:\n %v\n", string(output)), 0)
		}
	}

	// copy the replica binary, and configuration file to the nodes

	var wg sync.WaitGroup
	wg.Add(int(num_replicas_int))

	for j := int64(0); j < num_replicas_int; j++ {
		go func(i int) {
			nodes[i].Put_Load("protocols/ping/assets/ip_config.yaml", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/ping/assets/dummy", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			nodes[i].Put_Load("protocols/ping/assets/stats", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			wg.Done()
		}(int(j))
	}
	wg.Wait()
	ba.logger.Debug(fmt.Sprintf("Copied the ping binaries to all the nodes\n"), 0)

	return nil
}

func (ba *Ping) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
	replica_path := "/bench/dummy"
	stat_path := "/bench/stats"

	num_replicas, err := strconv.ParseInt(ba.options.Option["num_replicas"], 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}

	var wg sync.WaitGroup
	wg.Add(int(num_replicas))
	for i := 0; i < int(num_replicas); i++ {
		go func(j int) {
			nodes[j].ExecCmd("pkill -KILL -f dummy")
			nodes[j].ExecCmd("pkill -KILL -f stats")
			nodes[j].ExecCmd(fmt.Sprintf("rm -r %vbench/logs/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vbench/logs/", nodes[j].HomeDir))
			wg.Done()
		}(i)
	}
	wg.Wait()

	ba.logger.Debug(fmt.Sprintf("Killed all the replicas and clients\n"), 0)

	for j := 0; j < int(num_replicas); j++ {
		go func(i int) {
			nodes[i].ExecCmd("." + replica_path + " --name " + strconv.Itoa(i+1) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[i].HomeDir))
		}(j)
	}

	time.Sleep(15 * time.Second)

	ba.logger.Debug(fmt.Sprintf("Started all the replicas\n"), 0)

	go func() {
		nodes[0].ExecCmd("." + stat_path + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[0].HomeDir))
	}()

	ba.logger.Debug(fmt.Sprintf("Bootstrap complete\n"), 0)
	bootstrap_complete <- true

	time.Sleep(time.Duration(duration) * time.Second)

	var wg1 sync.WaitGroup
	wg1.Add(int(num_replicas))
	for j := 0; j < int(num_replicas); j++ {
		go func(i int) {
			nodes[i].ExecCmd("pkill -KILL -f dummy")
			nodes[i].ExecCmd("pkill -KILL -f stats")
			wg1.Done()
		}(j)
	}
	wg1.Wait()

	ba.logger.Debug(fmt.Sprintf("Killed all the replicas and stats client\n"), 0)

	result <- util.Performance{map[string]string{}}
}

func (ba *Ping) ExtractOptions(path string) protocols.ConsensusOptions {
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

	ba.logger.Debug(fmt.Sprintf("Ping options:\n %v\n", options.Option), 0)

	ba.options = options
	return options
}
