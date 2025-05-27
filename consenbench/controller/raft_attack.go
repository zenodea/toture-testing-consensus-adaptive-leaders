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
	fmt.Printf("Running RaftAttack1 continous leader node parition simplex  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {
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

type RaftAttack2 struct {
	logger *util.Logger
}

func NewRaftAttack2(logger *util.Logger) *RaftAttack2 {
	return &RaftAttack2{
		logger: logger,
	}
}

func (a *RaftAttack2) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running RaftAttack2 once leader node parition simplex  for %v seconds\n", duration)
	start_time := time.Now()

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

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {

	}

	fmt.Print("RaftAttack2 complete\n")
}

type RaftAttack3 struct {
	logger *util.Logger
}

func NewRaftAttack3(logger *util.Logger) *RaftAttack3 {
	return &RaftAttack3{
		logger: logger,
	}
}

func (a *RaftAttack3) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running RaftAttack3 leader is connected to a majority connected node for %v seconds\n", duration)
	start_time := time.Now()

	leaeder_node := oracle.GetTopNLeaders()[0]
	connected_node := (leaeder_node + 1) % len(nodes)

	fmt.Printf("attacking leader node %v\n", leaeder_node)

	for j := 0; j < len(nodes); j++ {
		if j == leaeder_node || j == connected_node {
			continue
		}
		links[leaeder_node][j].SetLoss(100)
		fmt.Printf("setting loss between %v and %v\n", leaeder_node, j)
	}

	fmt.Printf("attacking all nodes \n")

	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if i == j || i == leaeder_node || i == connected_node || j == connected_node {
				continue
			} else {
				links[i][j].SetLoss(100)
				fmt.Printf("setting loss between %v and %v\n", i, j)
			}
		}
	}

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("resetting attack\n")

	for j := 0; j < len(nodes); j++ {
		if j == leaeder_node || j == connected_node {
			continue
		}
		links[leaeder_node][j].SetLoss(0)
		fmt.Printf("setting loss between %v and %v\n", leaeder_node, j)
	}

	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if i == j || i == leaeder_node || i == connected_node || j == connected_node {
				continue
			} else {
				links[i][j].SetLoss(0)
				fmt.Printf("setting loss between %v and %v\n", i, j)
			}
		}
	}

	fmt.Print("RaftAttack3 complete\n")
}
