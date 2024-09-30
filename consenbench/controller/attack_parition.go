package controller

import (
	"fmt"
	"time"
	"toture-test/util"
)

type PartitionAttack struct {
	logger *util.Logger
}

func NewPartitionAttack(logger *util.Logger) *PartitionAttack {
	return &PartitionAttack{
		logger: logger,
	}
}

func (a *PartitionAttack) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack attack for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-5) {

	}

	fmt.Print("Partition attack complete\n")
}
