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

type Hotstuff_2 struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewHotstuff_2(logger *util.Logger) *Hotstuff_2 {
	return &Hotstuff_2{
		logger: logger,
	}
}

func (ba *Hotstuff_2) CopyConsensus(nodes []*common.Node) error {
	err := os.Chdir("protocols/hotstuff_2/assets/benchmark")
	if err != nil {
		log.Fatalf("Failed to change directory: %v", err)
	}

	// Define the command to run fab install
	cmd := exec.Command("fab", "install")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("Failed to run fab install: %v\n%v", err, string(output)))
	}
	return nil
}

func (ba *Hotstuff_2) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {

}

func (ba *Hotstuff_2) ExtractOptions(path string) protocols.ConsensusOptions {
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

	fmt.Printf("Hotstuff_2 options:\n %v\n", options.Option)

	ba.options = options
	return options
}

func (ba *Hotstuff_2) GetPerformance(outputs []string) util.Performance {
	return util.Performance{}
}
