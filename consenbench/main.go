package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"toture-test/consenbench/client"
	"toture-test/consenbench/controller"
)

func main() {
	// generic arguments
	is_controller := flag.Bool("is_controller", false, "is the node a controller")
	node_info_file := flag.String("node_config", "consenbench/assets/ip.yaml", "node ip configuration file")
	id := flag.Int("id", 1, "id of the node, according to node ip configuration file")
	debug_on := flag.Bool("debug_on", false, "turn on debug mode")
	debug_level := flag.Int("debug_level", 0, "debug level")
	logFilePath := flag.String("log_file_path", "bench/", "log file path")
	device := flag.String("device", "", "device to attack")

	// controller specific arguments
	attack_duration := flag.Int("attack_duration", 60, "duration of the attack in seconds")
	attack := flag.String("attack", "", "attack name to run")
	controller_operation_type := flag.String("controller_operation_type", "", "operation type of the controller: bootstrap/copy/run")
	consensus_algorithm := flag.String("consensus_algorithm", "raft", "consensus algorithm to run")

	flag.Parse()

	exit_timeout := time.Duration(10*(*attack_duration)) * time.Second
	exitTimer := time.AfterFunc(exit_timeout, func() {
		fmt.Println("Program exceeded maximum runtime, exiting forcefully.")
		os.Exit(1)
	})

	defer exitTimer.Stop()

	if *is_controller {
		options := controller.ControllerOptions{
			AttackDuration: *attack_duration,
			Attack:         *attack,
			NodeInfoFile:   *node_info_file,
			DebugOn:        *debug_on,
			DebugLevel:     *debug_level,
			LogFilePath:    *logFilePath,
			Device:         *device,
		}
		controller := controller.NewController(*id, options)
		if *controller_operation_type == "bootstrap" {
			controller.BootstrapClients()
		} else if *controller_operation_type == "copy" {
			controller.CopyConsensus(*consensus_algorithm)
		} else if *controller_operation_type == "run" {
			controller.Run(*consensus_algorithm)
		} else {
			panic("invalid controller operation type")
		}
	} else {
		options := client.ClientOptions{
			NodeInfoFile: *node_info_file,
			DebugOn:      *debug_on,
			DebugLevel:   *debug_level,
			LogFilePath:  *logFilePath,
			Device:       *device,
		}
		client := client.NewClient(*id, options)
		client.Run()
	}
}
