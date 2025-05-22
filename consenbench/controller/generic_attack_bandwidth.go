package controller

import (
	"fmt"
	"math/rand"
	"time"
	"toture-test/util"
)

// 1 egress links of the leader has low bandwidth

type BandwidthAttack_1 struct {
	logger *util.Logger
}

func NewBandwidthAttack_1(logger *util.Logger) *BandwidthAttack_1 {
	return &BandwidthAttack_1{
		logger: logger,
	}
}

func (a *BandwidthAttack_1) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running Bandwidth attack 1: changing the Bandwidth egress of leader for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		node_id := oracle.GetTopNLeaders()[0]
		// set all outgoing links to of leader to high Bandwidth
		fmt.Printf("attacking leader node %v\n", node_id)
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if i == node_id {
					links[i][j].SetBandwidth(100) // todo adjust according to the configured view change Bandwidth
					fmt.Printf("setting Bandwidth between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(1 * time.Second)
		fmt.Printf("resetting attack\n")
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if i == node_id {
					links[i][j].SetBandwidth(10000000)
					fmt.Printf("resetting Bandwidth between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print(" attack complete\n")
}

// 2 egress links of the 1st leader and second leader have low bandwidth

type BandwidthAttack_2 struct {
	logger *util.Logger
}

func NewBandwidthAttack_2(logger *util.Logger) *BandwidthAttack_2 {
	return &BandwidthAttack_2{
		logger: logger,
	}
}

func (a *BandwidthAttack_2) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running  Bandwidth attack 2 -- 1st and 2nd leader attacked for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		leader1 := oracle.GetTopNLeaders()[0]
		leader2 := oracle.GetTopNLeaders()[1]
		fmt.Printf("attacking leader nodes %v %v\n", leader1, leader2)
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if i == leader1 || i == leader2 {
					links[i][j].SetBandwidth(100) // todo adjust according to the configured view change Bandwidth
					fmt.Printf("setting Bandwidth between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(1 * time.Second)
		fmt.Printf("resetting attack\n")
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if i == leader1 || i == leader2 {
					links[i][j].SetBandwidth(10000000)
					fmt.Printf("resetting Bandwidth between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print(" attack complete\n")
}

// 3 egress links of the minority of the nodes have low bandwidth

type BandwidthAttack_3 struct {
	logger *util.Logger
}

func NewBandwidthAttack_3(logger *util.Logger) *BandwidthAttack_3 {
	return &BandwidthAttack_3{
		logger: logger,
	}
}

func (a *BandwidthAttack_3) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running  Bandwidth attack 3 -- minority nodes attacked for %v seconds\n", duration)
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
					links[i][j].SetBandwidth(100) // todo adjust according to the configured view change Bandwidth
					fmt.Printf("setting Bandwidth between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(1 * time.Second)
		fmt.Printf("resetting attack\n")
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if randNodes[i] {
					links[i][j].SetBandwidth(10000000)
					fmt.Printf("resetting Bandwidth between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print(" attack complete\n")
}

// 4 egress links of the majority of the nodes have low bandwidth

type BandwidthAttack_4 struct {
	logger *util.Logger
}

func NewBandwidthAttack_4(logger *util.Logger) *BandwidthAttack_4 {
	return &BandwidthAttack_4{
		logger: logger,
	}
}

func (a *BandwidthAttack_4) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running  Bandwidth attack 4 -- attacking majority of the nodes  for %v seconds\n", duration)
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
					links[i][j].SetBandwidth(100) // todo adjust according to the configured view change Bandwidth
					fmt.Printf("setting Bandwidth between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(1 * time.Second)
		fmt.Printf("resetting attack\n")
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if randNodes[i] {
					links[i][j].SetBandwidth(10000000)
					fmt.Printf("resetting Bandwidth between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print(" attack complete\n")
}
