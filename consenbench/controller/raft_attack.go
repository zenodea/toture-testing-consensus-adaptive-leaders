package controller

import (
	"fmt"
	"time"
	"toture-test/util"
)

type RaftAttack1 struct {
	logger *util.Logger
}

func NewRaftAttack1(logger *util.Logger) *RaftAttack1 {
	return &RaftAttack1{
		logger: logger,
	}
}

func (a *RaftAttack1) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-5) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		time.Sleep(1 * time.Second)
	}

	fmt.Print("RaftAttack1 complete\n")
}
