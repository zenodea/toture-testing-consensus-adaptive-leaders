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
	a.logger.Debug(fmt.Sprintf("Running Leader Parition Attack -- continous leader node parition duplex\n"), 0)

	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {
		a.logger.Debug(fmt.Sprintf("The leader order is %v\n", oracle.GetTopNLeaders()), 0)

		node_id := oracle.GetTopNLeaders()[0]

		a.logger.Debug(fmt.Sprintf("attacking leader node %v\n", node_id), 0)

		for j := 0; j < len(nodes); j++ {
			if node_id == j {
				continue
			}
			links[node_id][j].SetLoss(100)
			links[j][node_id].SetLoss(100)
			a.logger.Debug(fmt.Sprintf("setting duplex loss between %v and %v\n", node_id, j), 0)
		}

		time.Sleep(3 * time.Second)

		a.logger.Debug(fmt.Sprintf("resetting attack\n"), 0)

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

	a.logger.Debug(fmt.Sprintf("LeaderPartition attack complete\n"), 0)
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
	a.logger.Debug(fmt.Sprintf("Running One Quorum Parition Attack -- at a time only one quorum connected node\n"), 0)

	start_time := time.Now()
	good_node := 0

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {
		good_node = (good_node + 1) % len(nodes)

		a.logger.Debug(fmt.Sprintf("Only quorum connected node is %v\n", good_node), 0)

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j || i == good_node || j == good_node {
					continue
				}
				links[i][j].SetLoss(100)
				a.logger.Debug(fmt.Sprintf("setting loss between %v and %v\n", i, j), 0)
			}
		}

		time.Sleep(3 * time.Second)

		a.logger.Debug(fmt.Sprintf("resetting attack\n"), 0)

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

	a.logger.Debug(fmt.Sprintf("OneQuorumNodePartition complete\n"), 0)
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
	a.logger.Debug(fmt.Sprintf("Running Majority High Delay attack -- first majority nodes have high link delays\n"), 0)

	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {

		leaders := oracle.GetTopNLeaders()
		num_nodes := len(nodes)

		majority_leaders := make(map[int]bool)
		for i := 0; i < num_nodes/2+1; i++ {
			majority_leaders[leaders[i]] = true
		}

		a.logger.Debug(fmt.Sprintf("setting high delays for nodes %v\n", majority_leaders), 0)

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				if majority_leaders[i] && majority_leaders[j] {
					links[i][j].SetDelay(800)
					a.logger.Debug(fmt.Sprintf("setting delay between %v and %v\n", i, j), 0)
				}
			}
		}

		time.Sleep(3 * time.Second)

		a.logger.Debug(fmt.Sprintf("resetting attack\n"), 0)

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

	a.logger.Debug(fmt.Sprintf("MajorityHighDelay complete\n"), 0)
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
	a.logger.Debug(fmt.Sprintf("Running Minority Crash attack -- minority of the nodes crash\n"), 0)

	leaders := oracle.GetTopNLeaders()
	num_nodes := len(nodes)

	minority_leaders := make(map[int]bool)
	for i := 0; i < num_nodes/3; i++ {
		minority_leaders[leaders[i]] = true
	}

	a.logger.Debug(fmt.Sprintf("crashing minority nodes %v\n", minority_leaders), 0)

	for i := 0; i < len(nodes); i++ {
		if minority_leaders[i] {
			nodes[i].Pause()
			a.logger.Debug(fmt.Sprintf("crashed %v\n", i), 0)
		}
	}
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-5) {
		a.logger.Debug(fmt.Sprintf("The leader order is %v\n", oracle.GetTopNLeaders()), 0)
		time.Sleep(2 * time.Second)
	}

	for i := 0; i < len(nodes); i++ {
		if minority_leaders[i] {
			nodes[i].Continue()
		}
	}

	a.logger.Debug(fmt.Sprintf("MinorityCrash complete\n"), 0)
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
	a.logger.Debug(fmt.Sprintf("Running Minority Straggler -- minority nodes are slow\n"), 0)

	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {

		leaders := oracle.GetTopNLeaders()
		num_nodes := len(nodes)

		minority_leaders := make(map[int]bool)
		for i := 0; i < num_nodes/3; i++ {
			minority_leaders[leaders[i]] = true
		}

		a.logger.Debug(fmt.Sprintf("setting stragglers for nodes %v\n", minority_leaders), 0)

		for i := 0; i < len(nodes); i++ {
			if minority_leaders[i] {
				nodes[i].Pause()
				a.logger.Debug(fmt.Sprintf("paused %v\n", i), 0)
			}
		}

		time.Sleep(3 * time.Second)

		a.logger.Debug(fmt.Sprintf("resetting attack\n"), 0)

		for i := 0; i < len(nodes); i++ {
			if minority_leaders[i] {
				nodes[i].Continue()
			}
		}
	}

	time.Sleep(5 * time.Second)

	a.logger.Debug(fmt.Sprintf("MinorityStraggler complete\n"), 0)
}

type LeaderCrash struct {
	logger *util.Logger
}

func NewLeaderCrash(logger *util.Logger) *LeaderCrash {
	return &LeaderCrash{
		logger: logger,
	}
}

func (a *LeaderCrash) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	a.logger.Debug(fmt.Sprintf("Running Leader Crash attack\n"), 0)
	time.Sleep(30 * time.Second)
	leader := oracle.GetTopNLeaders()[0]
	nodes[leader].Kill()
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-35) {
		a.logger.Debug(fmt.Sprintf("The leader order is %v\n", oracle.GetTopNLeaders()), 0)
		time.Sleep(2 * time.Second)
	}

	a.logger.Debug(fmt.Sprintf("Leader Crash complete\n"), 0)
}
