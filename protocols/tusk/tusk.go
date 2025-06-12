package consensus

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"toture-test/consenbench/common"
	"toture-test/protocols"
	"toture-test/util"
)

type Tusk struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewTusk(logger *util.Logger) *Tusk {
	return &Tusk{
		logger: logger,
	}
}

func (ba *Tusk) CopyConsensus(nodes []*common.Node) error {
	println("Copying Tusk consensus to nodes using fabric")
	err := os.Chdir("protocols/tusk/assets/benchmark")
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
		fmt.Printf("fab install Output: %s\n", output)
	}
	return nil
}

func (ba *Tusk) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
	println("Running Tusk consensus using fabric")
	err := os.Chdir("protocols/tusk/assets/benchmark")
	if err != nil {
		panic("Failed to change directory")
	}
	// remove results directory and remake it
	cmd := exec.Command("rm", "-r", "results/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to delete results/ %v\n%v", err, string(output))
	} else {
		fmt.Printf("Deleted old results/ %s\n", output)
	}
	cmd = exec.Command("mkdir", "results/")
	output, err = cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("Failed to create results/ %v\n%v", err, string(output)))
	} else {
		fmt.Printf("Created  results/ %s\n", output)
	}

	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic(err.Error() + " while parsing num_replicas")

	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("error getting home directory:" + err.Error())
	}

	sshCmd := exec.Command("rm", []string{"-r", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		print("error while deleting logs/ " + err.Error() + " " + string(output) + "\n")
	} else {
		print("deleted local logs/ successfully\n" + string(output) + "\n")
	}

	sshCmd = exec.Command("mkdir", []string{"-p", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		panic("Error while creating logs/ " + err.Error() + " " + string(output) + "\n")
	} else {
		print("created logs/ successfully\n" + string(output) + "\n")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		cmd = exec.Command("fab", "remote", "--pid="+fmt.Sprintf("%v", os.Getpid()), "--attack-duration="+fmt.Sprintf("%v", duration), "--num-replicas="+num_replicas)
		output, err = cmd.CombinedOutput()
		if err != nil {
			panic(fmt.Sprintf("Failed to run %v: %v\n%v", cmd, err, string(output)))
		} else {
			// print output
			fmt.Printf("Fab install Output: %s\n", output)
			wg.Done()
		}
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// Notify the channel of specific signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine to handle signals
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println("Received signal from fabric:", sig)
		done <- true
	}()

	fmt.Println("Waiting for a signal from fabric (PID:", os.Getpid(), ")")
	<-done
	bootstrap_complete <- true
	fmt.Printf("bootstrap complete for tusk\n")
	wg.Wait()
	fmt.Printf("finished running tusk")
	p := ba.GetPerformance()
	err = os.Chdir(("../../../../"))
	if err != nil {
		panic("Failed to change directory")
	}
	result <- p
}

func (ba *Tusk) ExtractOptions(path string) protocols.ConsensusOptions {
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

	fmt.Printf("Tusk options:\n %v\n", options.Option)

	ba.options = options
	return options
}

func (ba *Tusk) GetPerformance() util.Performance {
	dirPath := "results/"

	// Read the directory to get the list of files
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	// Check if there is exactly one file in the directory
	if len(files) != 1 {
		log.Fatalf("Expected exactly one file in the directory, but found %d", len(files))
	}

	// Get the file name
	fileName := files[0].Name()

	// Full file path
	filePath := filepath.Join(dirPath, fileName)

	// Read the file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Convert the content to a string
	fileContent := string(content)

	return util.Performance{
		map[string]string{"summary": fileContent},
	}
}
