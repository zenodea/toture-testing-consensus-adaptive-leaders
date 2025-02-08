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
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		links[0][1].SetDelay(100)
		links[1][2].SetDelay(200)
		links[2][3].SetDelay(300)
		links[3][4].SetDelay(400)
		links[4][5].SetDelay(500)
		links[5][6].SetDelay(600)
		links[6][7].SetDelay(700)
		links[7][8].SetDelay(800)
		links[8][9].SetDelay(900)
		links[9][10].SetDelay(1000)
		links[10][11].SetDelay(1100)

		time.Sleep(3 * time.Second)

		links[0][1].SetDelay(0)
		links[1][2].SetDelay(0)
		links[2][3].SetDelay(0)
		links[3][4].SetDelay(0)
		links[4][5].SetDelay(0)
		links[5][6].SetDelay(0)
		links[6][7].SetDelay(0)
		links[7][8].SetDelay(0)
		links[8][9].SetDelay(0)
		links[9][10].SetDelay(0)
		links[10][11].SetDelay(0)

	}

	fmt.Print("Basic attack complete\n")
}
