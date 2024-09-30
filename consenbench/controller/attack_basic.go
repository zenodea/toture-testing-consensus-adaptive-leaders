package controller

import (
	"fmt"
	"time"
	"toture-test/util"
)

type BasicAttack struct {
	logger *util.Logger
}

func NewBasicAttack(logger *util.Logger) *BasicAttack {
	return &BasicAttack{
		logger: logger,
	}
}

func (a *BasicAttack) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running basic attack for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-5) {
		links[0][1].SetDelay(100)
		links[1][2].SetDelay(200)
		links[2][0].SetDelay(300)

		time.Sleep(1 * time.Second)

		links[0][1].SetDelay(0)
		links[1][2].SetDelay(0)
		links[2][0].SetDelay(0)

		time.Sleep(2 * time.Second)
	}

	fmt.Print("Basic attack complete\n")
}
