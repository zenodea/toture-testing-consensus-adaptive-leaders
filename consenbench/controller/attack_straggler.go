package controller

import (
	"fmt"
	"math/rand"
	"time"
	"toture-test/util"
)

// leader is a straggler

type StragglerAttack_1 struct {
	logger *util.Logger
}

func NewStragglerAttack_1(logger *util.Logger) *StragglerAttack_1 {
	return &StragglerAttack_1{
		logger: logger,
	}
}

func (a *StragglerAttack_1) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running straggler attack 1: straggler leader for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		node_id := oracle.GetTopNLeaders()[0]

		fmt.Printf("attacking leader node %v\n", node_id)

		nodes[node_id].Pause()
		time.Sleep(1 * time.Second)
		nodes[node_id].Continue()
		time.Sleep(3 * time.Second)
	}

	fmt.Print(" attack complete\n")
}

// leader 1 and 2 are stragglers

type StragglerAttack_2 struct {
	logger *util.Logger
}

func NewStragglerAttack_2(logger *util.Logger) *StragglerAttack_2 {
	return &StragglerAttack_2{
		logger: logger,
	}
}

func (a *StragglerAttack_2) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running straggler attack 2: straggler leader 1 and 2 for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		leader1 := oracle.GetTopNLeaders()[0]
		leader2 := oracle.GetTopNLeaders()[1]

		fmt.Printf("attacking leader nodes %v %v\n", leader1, leader2)

		nodes[leader1].Pause()
		nodes[leader2].Pause()
		time.Sleep(1 * time.Second)
		nodes[leader1].Continue()
		nodes[leader2].Continue()
		time.Sleep(3 * time.Second)
	}

	fmt.Print(" attack complete\n")
}

// minority stragglers

type StragglerAttack_3 struct {
	logger *util.Logger
}

func NewStragglerAttack_3(logger *util.Logger) *StragglerAttack_3 {
	return &StragglerAttack_3{
		logger: logger,
	}
}

func (a *StragglerAttack_3) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running straggler attack 3: minority stragglers for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		randNodes := make(map[int]bool)
		majority := len(nodes)/2 + 1
		minority := len(nodes) - majority

		for i := 0; i < minority; i++ {
			rand_node := rand.Intn(len(nodes))
			for randNodes[rand_node] {
				rand_node = rand.Intn(len(nodes))
			}
			randNodes[rand_node] = true
		}
		fmt.Printf("attacking minority nodes %v \n", randNodes)

		for node_id := range randNodes {
			nodes[node_id].Pause()
		}
		time.Sleep(1 * time.Second)
		for node_id := range randNodes {
			nodes[node_id].Continue()
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print(" attack complete\n")
}
