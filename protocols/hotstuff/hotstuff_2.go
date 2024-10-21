package consensus

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"toture-test/consenbench/common"
	"toture-test/protocols"
	"toture-test/util"
)

type Hotstuff struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewHotstuff(logger *util.Logger) *Hotstuff {
	return &Hotstuff{
		logger: logger,
	}
}

func (ba *Hotstuff) CopyConsensus(nodes []*common.Node) error {
	println("Copying Hotstuff consensus to nodes using fabric")
	err := os.Chdir("protocols/hotstuff_2/assets/benchmark")
	if err != nil {
		panic("Failed to change directory")
	}

	// Define the command to run fab install
	cmd := exec.Command("fab", "install")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("Failed to run fab install: %v\n%v", err, string(output)))
	} else {
		// print output
		fmt.Printf("Fab install Output: %s\n", output)
	}

	return nil
}

func (ba *Hotstuff) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {

}

func (ba *Hotstuff) ExtractOptions(path string) protocols.ConsensusOptions {
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

	fmt.Printf("Hotstuff options:\n %v\n", options.Option)

	ba.options = options
	return options
}

func (ba *Hotstuff) GetPerformance(outputs []string) util.Performance {
	return util.Performance{}
}
