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
		// set all incoming and outgoing links of rand_node to high loss
		fmt.Printf("attacking node %v\n", rand_node)
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if i == rand_node {
					links[i][j].SetLoss(100)
					fmt.Printf("setting loss between %v and %v\n", i, j)
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
				if i == rand_node {
					links[i][j].SetLoss(0)
					fmt.Printf("resetting loss between %v and %v\n", i, j)
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
	fmt.Printf("Running partition attack_2 randomly partition minority nodes  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		majority := len(nodes)/2 + 1
		minority := len(nodes) - majority

		rand_nodes := make(map[int]bool)

		for i := 0; i < minority; i++ {
			rand_node := rand.Intn(len(nodes))
			for rand_nodes[rand_node] {
				rand_node = rand.Intn(len(nodes))
			}
			rand_nodes[rand_node] = true
		}
		fmt.Printf("attacking nodes %v\n", rand_nodes)
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if rand_nodes[i] {
					links[i][j].SetLoss(100)
					fmt.Printf("setting loss between %v and %v\n", i, j)
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
				if rand_nodes[i] {
					links[i][j].SetLoss(0)
					fmt.Printf("resetting loss between %v and %v\n", i, j)
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
	fmt.Printf("Running partition attack_3 majority nodes attacked for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		rand_nodes := make(map[int]bool)

		for i := 0; i < len(nodes)/2+1; i++ {
			rand_node := rand.Intn(len(nodes))
			for rand_nodes[rand_node] {
				rand_node = rand.Intn(len(nodes))
			}
			rand_nodes[rand_node] = true

		}
		fmt.Printf("attacking nodes %v\n", rand_nodes)

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if rand_nodes[i] {
					links[i][j].SetLoss(100)
					fmt.Printf("setting loss between %v and %v\n", i, j)
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
				if rand_nodes[i] {
					links[i][j].SetLoss(0)
					fmt.Printf("resetting loss between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print("Partition attack 3 complete\n")
}

// 4 randomly partition the leader node in simplex

type PartitionAttack_4 struct {
	logger *util.Logger
}

func NewPartitionAttack_4(logger *util.Logger) *PartitionAttack_4 {
	return &PartitionAttack_4{
		logger: logger,
	}
}

func (a *PartitionAttack_4) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_4 leader node parition simplex  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		// select the leader
		node_id := oracle.GetTopNLeaders()[0]
		// set all outgoing links to of leader to high loss
		fmt.Printf("attacking leader node %v\n", node_id)
		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if i == node_id {
					links[i][j].SetLoss(100)
					fmt.Printf("setting loss between %v and %v\n", i, j)
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
					links[i][j].SetLoss(0)
					fmt.Printf("resetting loss between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print("Partition attack 4 complete\n")
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
	fmt.Printf("Running partition attack_5 -- each node can contact a minority for %v seconds\n", duration)
	start_time := time.Now()
	majority := len(nodes)/2 + 1
	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {

		for i := 0; i < len(nodes); i++ {
			fmt.Printf("attacking node %v\n", i)
			for j := i + 1; j < i+majority+1; j++ {
				links[i][j%len(nodes)].SetLoss(100)
				fmt.Printf("setting loss between %v and %v\n", i, j%len(nodes))
			}
		}

		time.Sleep(1 * time.Second)

		fmt.Printf("resetting attack\n")

		for i := 0; i < len(nodes); i++ {
			for j := i + 1; j < i+majority+1; j++ {
				links[i][j%len(nodes)].SetLoss(0)
				fmt.Printf("resetting loss between %v and %v\n", i, j%len(nodes))
			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print("Partition attack 5 complete\n")
}

// 6 randomly partition such that there is only one quorum connected node

type PartitionAttack_6 struct {
	logger *util.Logger
}

func NewPartitionAttack_6(logger *util.Logger) *PartitionAttack_6 {
	return &PartitionAttack_6{
		logger: logger,
	}
}

func (a *PartitionAttack_6) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_6 -- only one quorum connected node exists  for %v seconds\n", duration)
	start_time := time.Now()
	majority := len(nodes)/2 + 1
	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		// len(nodes) -  1 not attacked
		for i := 0; i < len(nodes)-1; i++ {
			fmt.Printf("attacking node %v\n", i)
			for j := i + 1; j < i+majority+1; j++ {
				if j%len(nodes) == len(nodes)-1 { // the node that we don't attack
					links[i][i-1].SetLoss(100) // because i -1 is not already chosen
					fmt.Printf("setting loss between %v and %v\n", i, i-1)
				} else {
					links[i][j%len(nodes)].SetLoss(100)
					fmt.Printf("setting loss between %v and %v\n", i, j%len(nodes))
				}
			}
		}

		time.Sleep(1 * time.Second)

		fmt.Printf("resetting attack\n")

		for i := 0; i < len(nodes)-1; i++ {
			for j := i + 1; j < i+majority+1; j++ {
				if j%len(nodes) == len(nodes)-1 {
					links[i][i-1].SetLoss(0)
					fmt.Printf("resetting loss between %v and %v\n", i, i-1)
				} else {
					links[i][j%len(nodes)].SetLoss(0)
					fmt.Printf("resetting loss between %v and %v\n", i, j%len(nodes))
				}

			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print("Partition 6 attack complete\n")
}

// 7 randomly partition such that there is a minority of quorum connected nodes

type PartitionAttack_7 struct {
	logger *util.Logger
}

func NewPartitionAttack_7(logger *util.Logger) *PartitionAttack_7 {
	return &PartitionAttack_7{
		logger: logger,
	}
}

func (a *PartitionAttack_7) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running partition attack_7 -- there is minority of quorum connected nodes  for %v seconds\n", duration)
	start_time := time.Now()
	majority := len(nodes)/2 + 1
	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
		for i := 0; i < majority; i++ {
			fmt.Printf("attacking node %v\n", i)
			for j := 0; j < majority; j++ {
				if i != j {
					links[i][j].SetLoss(100)
					fmt.Printf("setting loss between %v and %v\n", i, j)
				}
			}
		}

		time.Sleep(1 * time.Second)

		fmt.Printf("resetting attack\n")

		for i := 0; i < majority; i++ {
			for j := 0; j < majority; j++ {
				if i != j {
					links[i][j].SetLoss(0)
					fmt.Printf("resetting loss between %v and %v\n", i, j)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}

	fmt.Print("Partition 7 attack complete\n")
}
