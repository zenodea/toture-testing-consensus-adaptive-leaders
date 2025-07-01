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
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"toture-test/consenbench/common"
	"toture-test/protocols"
	"toture-test/util"
)

type Bullshark struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewBullshark(logger *util.Logger) *Bullshark {
	return &Bullshark{
		logger: logger,
	}
}

func (ba *Bullshark) CopyConsensus(nodes []*common.Node) error {
	ba.logger.Debug(fmt.Sprintf("Copying Bullshark consensus to nodes using fabric"), 0)
	err := os.Chdir("protocols/bullshark/assets/benchmark")
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
		ba.logger.Debug(fmt.Sprintf("fab install Output: %s\n", output), 0)
	}
	return nil
}

func (ba *Bullshark) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool) {
	ba.logger.Debug(fmt.Sprintf("Running Bullshark consensus using fabric"), 0)
	err := os.Chdir("protocols/bullshark/assets/benchmark")
	if err != nil {
		panic("Failed to change directory")
	}
	// remove results directory and remake it
	cmd := exec.Command("rm", "-r", "results/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug(fmt.Sprintf("Failed to delete results/ %v\n%v", err, string(output)), 0)
	} else {
		ba.logger.Debug(fmt.Sprintf("Deleted old results/ %s\n", output), 0)
	}
	cmd = exec.Command("mkdir", "results/")
	output, err = cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("Failed to create results/ %v\n%v", err, string(output)))
	} else {
		ba.logger.Debug(fmt.Sprintf("Created  results/ %s\n", output), 0)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("error getting home directory:" + err.Error())
	}

	sshCmd := exec.Command("rm", []string{"-r", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug(fmt.Sprintf("error while deleting logs/ "+err.Error()+" "+string(output)+"\n"), 0)
	} else {
		ba.logger.Debug(fmt.Sprintf("deleted local logs/ successfully\n"+string(output)+"\n"), 0)
	}

	sshCmd = exec.Command("mkdir", []string{"-p", filepath.Join(homeDir, "toture-testing-consensus/logs")}...)
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		panic("Error while creating logs/ " + err.Error() + " " + string(output) + "\n")
	} else {
		ba.logger.Debug(fmt.Sprintf("created logs/ successfully\n"+string(output)+"\n"), 0)
	}

	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic(err.Error() + " while parsing num_replicas")

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
			ba.logger.Debug(fmt.Sprintf("Fab remote Output: %s\n", output), 0)
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

		ba.logger.Debug(fmt.Sprintf("Received signal from fabric:", sig), 0)
		done <- true
	}()

	ba.logger.Debug(fmt.Sprintf("Waiting for a signal from fabric (PID:", os.Getpid(), ")"), 0)
	<-done
	time.Sleep(10 * time.Second)
	bootstrap_complete <- true
	ba.logger.Debug(fmt.Sprintf("bootstrap complete for bullshark\n"), 0)
	wg.Wait()
	ba.logger.Debug(fmt.Sprintf("finished running bullshark"), 0)
	p := ba.GetPerformance(duration)
	err = os.Chdir(("../../../../"))
	if err != nil {
		panic("Failed to change directory")
	}

	cmd = exec.Command("pkill", "fab")
	output, err = sshCmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug(fmt.Sprintf("Error while killing fab "+err.Error()+" "+string(output)+"\n"), 0)
	}
	result <- p
}

func (ba *Bullshark) ExtractOptions(path string) protocols.ConsensusOptions {
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

	ba.logger.Debug(fmt.Sprintf("Bullshark options:\n %v\n", options.Option), 0)

	ba.options = options
	return options
}

func (ba *Bullshark) GetPerformance(duration int) util.Performance {
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

	p := util.Performance{
		map[string]string{"summary": fileContent},
	}

	var throughput string
	var latency string
	var run_duration string

	lines := strings.Split(fileContent, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Consensus TPS") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				valuePart := strings.TrimSpace(parts[1])
				tokens := strings.Split(valuePart, " ")
				if len(tokens) > 0 {
					throughput = strings.ReplaceAll(tokens[0], ",", "")
				}
			}
		} else if strings.HasPrefix(line, "Consensus latency") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				valuePart := strings.TrimSpace(parts[1])
				tokens := strings.Split(valuePart, " ")
				if len(tokens) > 0 {
					latency = strings.ReplaceAll(tokens[0], ",", "")
				}
			}
		} else if strings.HasPrefix(line, "Execution time") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				valuePart := strings.TrimSpace(parts[1])
				tokens := strings.Split(valuePart, " ")
				if len(tokens) > 0 {
					run_duration = strings.ReplaceAll(tokens[0], ",", "")
				}
			}
		}
	}

	throughput_int, _ := strconv.Atoi(throughput)
	duration_int, _ := strconv.Atoi(run_duration)

	throuhgput_exact := (throughput_int * duration_int) / duration

	fmt.Printf("%v,%v,%v,", throuhgput_exact, latency, 0)

	return p
}
