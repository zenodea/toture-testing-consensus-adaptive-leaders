package controller

import (
	"fmt"
	"time"
	"toture-test/util"
)

type LeaderPartition struct {
	logger *util.Logger
}

func NewLeaderPartition(logger *util.Logger) *LeaderPartition {
	return &LeaderPartition{
		logger: logger,
	}
}

func (a *LeaderPartition) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running Leader Parition Attack -- continous leader node parition duplex\n")

	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())

		node_id := oracle.GetTopNLeaders()[0]

		fmt.Printf("attacking leader node %v\n", node_id)

		for j := 0; j < len(nodes); j++ {
			if node_id == j {
				continue
			}
			links[node_id][j].SetLoss(100)
			links[j][node_id].SetLoss(100)
			fmt.Printf("setting duplex loss between %v and %v\n", node_id, j)
		}

		time.Sleep(3 * time.Second)

		fmt.Printf("resetting attack\n")

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				links[i][j].SetLoss(0)
			}
		}
	}

	time.Sleep(5 * time.Second)

	fmt.Print("LeaderPartition attack complete\n")
}

type OneQuorumNodePartition struct {
	logger *util.Logger
}

func NewOneQuorumNodePartition(logger *util.Logger) *OneQuorumNodePartition {
	return &OneQuorumNodePartition{
		logger: logger,
	}
}

func (a *OneQuorumNodePartition) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running One Quorum Parition Attack -- at a time only one quorum connected node\n")

	start_time := time.Now()
	good_node := 0

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {
		good_node = (good_node + 1) % len(nodes)

		fmt.Printf("Only quorum connected node is %v\n", good_node)

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j || i == good_node || j == good_node {
					continue
				}
				links[i][j].SetLoss(100)
				fmt.Printf("setting loss between %v and %v\n", i, j)
			}
		}

		time.Sleep(3 * time.Second)

		fmt.Printf("resetting attack\n")

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				links[i][j].SetLoss(0)
			}
		}
	}

	time.Sleep(5 * time.Second)

	fmt.Print("OneQuorumNodePartition complete\n")
}

type MajorityHighDelay struct {
	logger *util.Logger
}

func NewMajorityHighDelay(logger *util.Logger) *MajorityHighDelay {
	return &MajorityHighDelay{
		logger: logger,
	}
}

func (a *MajorityHighDelay) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running Majority High Delay attack -- first majority nodes have high link delays\n")

	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {

		leaders := oracle.GetTopNLeaders()
		num_nodes := len(nodes)

		majority_leaders := make(map[int]bool)
		for i := 0; i < num_nodes/2+1; i++ {
			majority_leaders[leaders[i]] = true
		}

		fmt.Printf("setting high delays for nodes %v\n", majority_leaders)

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if majority_leaders[i] && majority_leaders[j] {
					links[i][j].SetDelay(800)
					fmt.Printf("setting delay between %v and %v\n", i, j)
				}
			}
		}

		time.Sleep(3 * time.Second)

		fmt.Printf("resetting attack\n")

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				links[i][j].SetDelay(0)
			}
		}
	}

	time.Sleep(5 * time.Second)

	fmt.Print("MajorityHighDelay complete\n")
}

type MinorityCrash struct {
	logger *util.Logger
}

func NewMinorityCrash(logger *util.Logger) *MinorityCrash {
	return &MinorityCrash{
		logger: logger,
	}
}

func (a *MinorityCrash) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running Minority Crash attack -- minority of the nodes crash\n")

	leaders := oracle.GetTopNLeaders()
	num_nodes := len(nodes)

	minority_leaders := make(map[int]bool)
	for i := 0; i < num_nodes/3; i++ {
		minority_leaders[leaders[i]] = true
	}

	fmt.Printf("crashing minority nodes %v\n", minority_leaders)

	for i := 0; i < len(nodes); i++ {
		if minority_leaders[i] {
			nodes[i].Pause()
			fmt.Printf("crashed %v\n", i)
		}
	}
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-5) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		time.Sleep(2 * time.Second)
	}

	for i := 0; i < len(nodes); i++ {
		if minority_leaders[i] {
			nodes[i].Continue()
		}
	}

	fmt.Print("MinorityCrash complete\n")
}

type MinorityStraggler struct {
	logger *util.Logger
}

func NewMinorityStraggler(logger *util.Logger) *MinorityStraggler {
	return &MinorityStraggler{
		logger: logger,
	}
}

func (a *MinorityStraggler) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running Minority Straggler -- minority nodes are slow\n")

	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {

		leaders := oracle.GetTopNLeaders()
		num_nodes := len(nodes)

		minority_leaders := make(map[int]bool)
		for i := 0; i < num_nodes/3; i++ {
			minority_leaders[leaders[i]] = true
		}

		fmt.Printf("setting stragglers for nodes %v\n", minority_leaders)

		for i := 0; i < len(nodes); i++ {
			if minority_leaders[i] {
				nodes[i].Pause()
				fmt.Printf("paused %v\n", i)
			}
		}

		time.Sleep(3 * time.Second)

		fmt.Printf("resetting attack\n")

		for i := 0; i < len(nodes); i++ {
			if minority_leaders[i] {
				nodes[i].Continue()
			}
		}
	}

	time.Sleep(5 * time.Second)

	fmt.Print("MinorityStraggler complete\n")
}
