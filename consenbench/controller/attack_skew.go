package controller

import (
	"fmt"
	"math/rand"
	"time"
	"toture-test/util"
)

// SkewAttack_1 is an attack that changes the skew of the leader node to all other nodes

type SkewAttack_1 struct {
	logger *util.Logger
}

func NewSkewAttack_1(logger *util.Logger) *SkewAttack_1 {
	return &SkewAttack_1{
		logger: logger,
	}
}

func (a *SkewAttack_1) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running Skew Attack_1 : changing the skew of leader for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		node_id := oracle.GetTopNLeaders()[0]
		// set all outgoing links to of leader to high delay
		fmt.Printf("attacking leader node %v\n", node_id)
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if i == node_id {
					links[i][j].SetDelay(10)
					links[i][j].SetDelay(20)
					links[i][j].SetDelay(30)
					links[i][j].SetDelay(40)
					links[i][j].SetDelay(50)
					links[i][j].SetDelay(0)
					fmt.Printf("setting skew between %v and %v\n", i, j)
				}
			}
		}

		time.Sleep(3 * time.Second)
	}

	fmt.Print(" attack complete\n")
}

// SkewAttack_2 is an attack that changes the skew of the minority node to all other nodes

type SkewAttack_2 struct {
	logger *util.Logger
}

func NewSkewAttack_2(logger *util.Logger) *SkewAttack_2 {
	return &SkewAttack_2{
		logger: logger,
	}
}

func (a *SkewAttack_2) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running Skew Attack_2 : changing the skew of minority for %v seconds\n", duration)
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
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if randNodes[i] {
					links[i][j].SetDelay(10) // todo adjust according to the configured view change delay
					links[i][j].SetDelay(20) // todo adjust according to the configured view change delay
					links[i][j].SetDelay(30) // todo adjust according to the configured view change delay
					links[i][j].SetDelay(40) // todo adjust according to the configured view change delay
					links[i][j].SetDelay(50) // todo adjust according to the configured view change delay
					links[i][j].SetDelay(0)  // todo adjust according to the configured view change delay
					fmt.Printf("setting skew between %v and %v\n", i, j)
				}
			}
		}

		time.Sleep(3 * time.Second)
	}

	fmt.Print(" attack complete\n")
}

// SkewAttack_3 is an attack that changes the skew of the majority node to all other nodes

type SkewAttack_3 struct {
	logger *util.Logger
}

func NewSkewAttack_3(logger *util.Logger) *SkewAttack_3 {
	return &SkewAttack_3{
		logger: logger,
	}
}

func (a *SkewAttack_3) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running Skew Attack_3 : changing the skew of majorty for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		randNodes := make(map[int]bool)

		majority := len(nodes)/2 + 1

		for i := 0; i < majority; i++ {
			rand_node := rand.Intn(len(nodes))
			for randNodes[rand_node] {
				rand_node = rand.Intn(len(nodes))
			}
			randNodes[rand_node] = true
		}
		fmt.Printf("attacking majority nodes %v \n", randNodes)
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if randNodes[i] {
					links[i][j].SetDelay(10) // todo adjust according to the configured view change delay
					links[i][j].SetDelay(20) // todo adjust according to the configured view change delay
					links[i][j].SetDelay(30) // todo adjust according to the configured view change delay
					links[i][j].SetDelay(40) // todo adjust according to the configured view change delay
					links[i][j].SetDelay(50) // todo adjust according to the configured view change delay
					links[i][j].SetDelay(0)  // todo adjust according to the configured view change delay
					fmt.Printf("setting skew between %v and %v\n", i, j)
				}
			}
		}

		time.Sleep(3 * time.Second)
	}

	fmt.Print(" attack complete\n")
}
