package controller

import (
	"fmt"
	"time"
	"toture-test/util"
)

type LeaderAttack1 struct {
	logger *util.Logger
}

func NewLeaderAttack1(logger *util.Logger) *LeaderAttack1 {
	return &LeaderAttack1{
		logger: logger,
	}
}

func (a *LeaderAttack1) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running LeaderAttack1 continous leader node parition simplex  for %v seconds\n", duration)
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
		time.Sleep(3 * time.Second)

		fmt.Printf("resetting attack\n")

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				}
				links[i][j].SetLoss(0)
				fmt.Printf("resetting loss between %v and %v\n", i, j)
			}

		}
	}

	fmt.Print("LeaderAttack1 complete\n")
}

type LeaderAttack2 struct {
	logger *util.Logger
}

func NewLeaderAttack2(logger *util.Logger) *LeaderAttack2 {
	return &LeaderAttack2{
		logger: logger,
	}
}

func (a *LeaderAttack2) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running LeaderAttack2 once leader node parition simplex  for %v seconds\n", duration)
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
	time.Sleep(3 * time.Second)

	fmt.Printf("resetting attack\n")

	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if i == j {
				continue
			}
			links[i][j].SetLoss(0)
			fmt.Printf("resetting loss between %v and %v\n", i, j)
		}

	}

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {

	}

	fmt.Print("LeaderAttack2 complete\n")
}

type LeaderAttack3 struct {
	logger *util.Logger
}

func NewLeaderAttack3(logger *util.Logger) *LeaderAttack3 {
	return &LeaderAttack3{
		logger: logger,
	}
}

func (a *LeaderAttack3) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running LeaderAttack3 leader is connected to a majority connected node for %v seconds\n", duration)
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

	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if i == j {
				continue
			} else {
				links[i][j].SetLoss(0)
				fmt.Printf("setting loss between %v and %v\n", i, j)
			}
		}
	}

	fmt.Print("LeaderAttack3 complete\n")
}

type LeaderAttack4 struct {
	logger *util.Logger
}

func NewLeaderAttack4(logger *util.Logger) *LeaderAttack4 {
	return &LeaderAttack4{
		logger: logger,
	}
}

func (a *LeaderAttack4) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {

	fmt.Printf("Running LeaderAttack4 in each epoch only one quorum connected node for %v seconds\n", duration)

	start_time := time.Now()
	good_node := 0

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {

		good_node = (good_node + 1) % len(nodes)

		fmt.Printf("attacking all nodes except %v \n", good_node)

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j || i == good_node || j == good_node {
					continue
				} else {
					links[i][j].SetLoss(100)
					fmt.Printf("setting loss between %v and %v\n", i, j)
				}
			}
		}

		time.Sleep(3 * time.Second)

		fmt.Printf("resetting attack\n")

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j {
					continue
				} else {
					links[i][j].SetLoss(0)
					fmt.Printf("setting loss between %v and %v\n", i, j)
				}
			}
		}
	}
	fmt.Print("LeaderAttack4 complete\n")
}

type LeaderAttack5 struct {
	logger *util.Logger
}

func NewLeaderAttack5(logger *util.Logger) *LeaderAttack5 {
	return &LeaderAttack5{
		logger: logger,
	}
}

func (a *LeaderAttack5) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {

	fmt.Printf("Running LeaderAttack5 leader high delay but no view changes for %v seconds\n", duration)

	start_time := time.Now()

	leaeder_node := oracle.GetTopNLeaders()[0]

	fmt.Printf("attacking leader node %v\n", leaeder_node)

	for j := 0; j < len(nodes); j++ {
		if j == leaeder_node {
			continue
		}
		links[leaeder_node][j].SetDelay(20)
		fmt.Printf("setting delay between %v and %v\n", leaeder_node, j)
	}

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		time.Sleep(500 * time.Millisecond)
	}

	for j := 0; j < len(nodes); j++ {
		if j == leaeder_node {
			continue
		}
		links[leaeder_node][j].SetDelay(0)
		fmt.Printf("setting delay between %v and %v\n", leaeder_node, j)
	}

	fmt.Print("LeaderAttack5 complete\n")
}

type LeaderAttack6 struct {
	logger *util.Logger
}

func NewLeaderAttack6(logger *util.Logger) *LeaderAttack6 {
	return &LeaderAttack6{
		logger: logger,
	}
}

func (a *LeaderAttack6) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {

	fmt.Printf("Running LeaderAttack6 all links have close to view timeout delay for %v seconds\n", duration)

	start_time := time.Now()

	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if i == j {
				continue
			}
			links[i][j].SetDelay(300)
		}
	}

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		time.Sleep(500 * time.Millisecond)
	}

	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if i == j {
				continue
			}
			links[i][j].SetDelay(0)
		}
	}

	fmt.Print("LeaderAttack6 complete\n")
}

type LeaderAttack7 struct {
	logger *util.Logger
}

func NewLeaderAttack7(logger *util.Logger) *LeaderAttack7 {
	return &LeaderAttack7{
		logger: logger,
	}
}

func (a *LeaderAttack7) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {

	fmt.Printf("Running LeaderAttack7 two stragglers for %v seconds\n", duration)

	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		nodes[0].Pause()
		nodes[1].Pause()
		time.Sleep(3 * time.Second)
		nodes[0].Continue()
		nodes[1].Continue()
	}

	fmt.Print("LeaderAttack7 complete\n")
}

type LeaderAttack8 struct {
	logger *util.Logger
}

func NewLeaderAttack8(logger *util.Logger) *LeaderAttack8 {
	return &LeaderAttack8{
		logger: logger,
	}
}

func (a *LeaderAttack8) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {

	fmt.Printf("Running LeaderAttack8 straggler in closest majority for %v seconds\n", duration)

	start_time := time.Now()

	n := len(nodes)
	m := n/2 + 1

	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if i == j {
				continue
			}
			if i < m && j < m {
				continue
			}
			if i >= m && j >= m {
				continue
			}
			links[i][j].SetDelay(10)
		}
	}

	fmt.Printf("majority and minority set\n")

	// run for just 10 seconds

	for time.Now().Sub(start_time).Seconds() < float64(10) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("Setting node 0 straggler\n")

	for time.Now().Sub(start_time).Seconds() < float64(duration-20) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		nodes[0].Pause() // straggler in majority
		time.Sleep(1 * time.Second)
		nodes[0].Continue() // straggler in majority
	}

	nodes[0].Continue()

	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if i == j {
				continue
			}
			links[i][j].SetDelay(0)
		}
	}

	fmt.Print("LeaderAttack8 complete\n")
}

type LeaderAttack9 struct {
	logger *util.Logger
}

func NewLeaderAttack9(logger *util.Logger) *LeaderAttack9 {
	return &LeaderAttack9{
		logger: logger,
	}
}

func (a *LeaderAttack9) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {

	fmt.Printf("Running LeaderAttack9 leader crash for %v seconds\n", duration)

	start_time := time.Now()

	// run for just 10 seconds

	for time.Now().Sub(start_time).Seconds() < float64(10) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("killing leader\n")

	nodes[oracle.GetTopNLeaders()[0]].Kill()

	for time.Now().Sub(start_time).Seconds() < float64(duration-20) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Print("LeaderAttack9 complete\n")
}
