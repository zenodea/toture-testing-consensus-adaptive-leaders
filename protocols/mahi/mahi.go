package consensus

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
	"toture-test/consenbench/common"
	"toture-test/protocols"
	"toture-test/util"
)

type Mahi struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewMahi(logger *util.Logger) *Mahi {
	return &Mahi{
		logger: logger,
	}
}

func (ba *Mahi) CopyConsensus(nodes []*common.Node) error {
	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("num_replicas not found in options")
	}

	num_replicas_int, err := strconv.ParseInt(num_replicas, 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}
	if num_replicas_int > int64(len(nodes)) {
		panic("Not enough nodes to deploy mahi")
	}

	nodes[0].ExecCmd(fmt.Sprintf("sudo rm -r async-mystecity; git clone https://github.com/PasinduTennage/async-mystecity; cd async-mystecity; git checkout consensus-rework; sudo apt-get install -y libfontconfig1-dev; source %v.cargo/env; cargo build", nodes[0].HomeDir))

	println("Cloned the mahi repository and built the binary")

	nodes[0].Get_Load(fmt.Sprintf("%vasync-mystecity/target/debug/mysticeti", nodes[0].HomeDir), "protocols/mahi/assets/")

	println("Copied the mahi binary to the local machine")

	var wg sync.WaitGroup
	wg.Add(int(num_replicas_int))

	for j := int64(0); j < num_replicas_int; j++ {
		go func(i int) {
			nodes[i].Put_Load("protocols/mahi/assets/mysticeti", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
			wg.Done()
		}(int(j))
	}
	wg.Wait()
	fmt.Print("Copied the mahi binary to all the nodes\n")
	return nil
}

func (ba *Mahi) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {

}

func (ba *Mahi) ExtractOptions(path string) protocols.ConsensusOptions {
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

	fmt.Printf("Mahi Mahi options:\n %v\n", options.Option)

	ba.options = options
	return options
}

func (ba *Mahi) GetPerformance(outputs []string) util.Performance {
	return util.Performance{}
}
