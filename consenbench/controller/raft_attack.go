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
	fmt.Printf("Running RaftAttack1 leader node parition simplex  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		// select the leader
		node_id := oracle.GetTopNLeaders()[0]
		// set all outgoing links to of leader to high loss
		fmt.Printf("attacking leader node %v\n", node_id)
		for i := 0; i < len(nodes); i++ {
			if i == node_id {
				for j := 0; j < len(nodes); j++ {
					if i == j {
						continue
					}
					links[i][j].SetLoss(100)
					fmt.Printf("setting loss between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(2 * time.Second)

		fmt.Printf("resetting attack\n")

		for i := 0; i < len(nodes); i++ {
			if i == node_id {
				for j := 0; j < len(nodes); j++ {
					if i == j {
						continue
					}
					links[i][j].SetLoss(0)
					fmt.Printf("resetting loss between %v and %v\n", i, j)
				}
			}
		}
	}

	fmt.Print("RaftAttack1 complete\n")
}
