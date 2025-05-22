package controller

import (
	"fmt"
	"time"
	"toture-test/util"
)

type NoopAttack struct {
	logger *util.Logger
}

func NewNoopAttack(logger *util.Logger) *NoopAttack {
	return &NoopAttack{
		logger: logger,
	}
}

func (a *NoopAttack) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-5) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		time.Sleep(1 * time.Second)
	}

	fmt.Print("Noop attack complete\n")
}
