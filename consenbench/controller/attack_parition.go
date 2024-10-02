package controller

import (
	"fmt"
	"math/rand"
	"time"
	"toture-test/util"
)

// 1 randomly partition one node at a time

type PartitionAttack_1 struct {
	logger *util.Logger
}

func NewPartitionAttack_1(logger *util.Logger) *PartitionAttack_1 {
	return &PartitionAttack_1{
		logger: logger,
	}
}

func (a *PartitionAttack_1) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_1:randomly partition one node at a time  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		// select a random node
		rand_node := rand.Intn(len(nodes))
		// set all incoming and outgoing links to of rand_node to high loss
		fmt.Printf("attacking node %v\n", rand_node+2)
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if i == rand_node || j == rand_node {
					links[i][j].SetLoss(100)
					fmt.Printf("setting loss between %v and %v\n", i+2, j+2)
				}
			}
		}
		time.Sleep(1 * time.Second)
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if i == rand_node || j == rand_node {
					links[i][j].SetLoss(0)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print("Partition attack 1 complete\n")
}

// 2 randomly partition minority node at a time

type PartitionAttack_2 struct {
	logger *util.Logger
}

func NewPartitionAttack_2(logger *util.Logger) *PartitionAttack_2 {
	return &PartitionAttack_2{
		logger: logger,
	}
}

func (a *PartitionAttack_2) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_2  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		rand_nodes := make(map[int]bool)
		rand_nodes_list := make([]int, len(nodes)/2-1)
		for i := 0; i < len(nodes)/2-1; i++ {
			rand_node := rand.Intn(len(nodes))
			for rand_nodes[rand_node] {
				rand_node = rand.Intn(len(nodes))
			}
			rand_nodes[rand_node] = true
			rand_nodes_list[i] = rand_node + 2
		}
		fmt.Printf("attacking node %v\n", rand_nodes_list)
		// set all incoming and outgoing links to of rand_nodes to high loss
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if rand_nodes[i] || rand_nodes[j] {
					links[i][j].SetLoss(100)
					fmt.Printf("setting loss between %v and %v\n", i+2, j+2)
				}
			}
		}
		time.Sleep(1 * time.Second)
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if rand_nodes[i] || rand_nodes[j] {
					links[i][j].SetLoss(0)
					fmt.Printf("resetting loss between %v and %v\n", i+2, j+2)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print("Partition attack 2 attack complete\n")
}

// 3 randomly partition majority node at a time

type PartitionAttack_3 struct {
	logger *util.Logger
}

func NewPartitionAttack_3(logger *util.Logger) *PartitionAttack_3 {
	return &PartitionAttack_3{
		logger: logger,
	}
}

func (a *PartitionAttack_3) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_3  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		rand_nodes := make(map[int]bool)
		rand_nodes_list := make([]int, len(nodes)/2+1)
		for i := 0; i < len(nodes)/2+1; i++ {
			rand_node := rand.Intn(len(nodes))
			for rand_nodes[rand_node] {
				rand_node = rand.Intn(len(nodes))
			}
			rand_nodes[rand_node] = true
			rand_nodes_list[i] = rand_node + 2
		}
		fmt.Printf("attacking nodes %v\n", rand_nodes_list)
		// set all incoming and outgoing links to of rand_nodes to high loss
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if rand_nodes[i] || rand_nodes[j] {
					links[i][j].SetLoss(100)
					fmt.Printf("setting loss between %v and %v\n", i+2, j+2)
				}
			}
		}
		time.Sleep(1 * time.Second)
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if rand_nodes[i] || rand_nodes[j] {
					links[i][j].SetLoss(0)
					fmt.Printf("resetting loss between %v and %v\n", i+2, j+2)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print("Partition attack 3 complete\n")
}

// 4_1 randomly partition the leader node in simplex

type PartitionAttack_4_1 struct {
	logger *util.Logger
}

func NewPartitionAttack_4_1(logger *util.Logger) *PartitionAttack_4_1 {
	return &PartitionAttack_4_1{
		logger: logger,
	}
}

func (a *PartitionAttack_4_1) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_4_1  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {

	}

	fmt.Print("Partition attack complete\n")
}

// 4_2 randomly partition the leader node in duplex

type PartitionAttack_4_2 struct {
	logger *util.Logger
}

func NewPartitionAttack_4_2(logger *util.Logger) *PartitionAttack_4_2 {
	return &PartitionAttack_4_2{
		logger: logger,
	}
}

func (a *PartitionAttack_4_2) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_4_2  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {

	}

	fmt.Print("Partition attack complete\n")
}

// 5 randomly partition such that each node can contact only a minority of other nodes

type PartitionAttack_5 struct {
	logger *util.Logger
}

func NewPartitionAttack_5(logger *util.Logger) *PartitionAttack_5 {
	return &PartitionAttack_5{
		logger: logger,
	}
}

func (a *PartitionAttack_5) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_5  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {

	}

	fmt.Print("Partition attack complete\n")
}

// 6 randomly partition such that there is only on quorum connected node

type PartitionAttack_6 struct {
	logger *util.Logger
}

func NewPartitionAttack_6(logger *util.Logger) *PartitionAttack_6 {
	return &PartitionAttack_6{
		logger: logger,
	}
}

func (a *PartitionAttack_6) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_6  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {

	}

	fmt.Print("Partition attack complete\n")
}

// 7 randomly partial connectivity where links are up and down randomly

type PartitionAttack_7 struct {
	logger *util.Logger
}

func NewPartitionAttack_7(logger *util.Logger) *PartitionAttack_7 {
	return &PartitionAttack_7{
		logger: logger,
	}
}

func (a *PartitionAttack_7) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_7  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {

	}

	fmt.Print("Partition attack complete\n")
}

// 8 randomly partition such that there is a miniroty of quorum connected nodes

type PartitionAttack_8 struct {
	logger *util.Logger
}

func NewPartitionAttack_8(logger *util.Logger) *PartitionAttack_8 {
	return &PartitionAttack_8{
		logger: logger,
	}
}

func (a *PartitionAttack_8) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_8  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {

	}

	fmt.Print("Partition attack complete\n")
}
