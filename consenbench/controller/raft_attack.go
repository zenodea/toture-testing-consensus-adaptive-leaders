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

type RaftAttack4 struct {
	logger *util.Logger
}

func NewRaftAttack4(logger *util.Logger) *RaftAttack4 {
	return &RaftAttack4{
		logger: logger,
	}
}

func (a *RaftAttack4) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {

	fmt.Printf("Running RaftAttack4 in each epoch only one quorum connected node for %v seconds\n", duration)

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

		time.Sleep(2 * time.Second)

		fmt.Printf("resetting attack\n")

		for i := 0; i < len(nodes); i++ {
			for j := 0; j < len(nodes); j++ {
				if i == j || i == good_node || j == good_node {
					continue
				} else {
					links[i][j].SetLoss(0)
					fmt.Printf("setting loss between %v and %v\n", i, j)
				}
			}
		}
	}
	fmt.Print("RaftAttack4 complete\n")
}

type RaftAttack5 struct {
	logger *util.Logger
}

func NewRaftAttack5(logger *util.Logger) *RaftAttack5 {
	return &RaftAttack5{
		logger: logger,
	}
}

func (a *RaftAttack5) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {

	fmt.Printf("Running RaftAttack5 leader high delay but no view changes for %v seconds\n", duration)

	start_time := time.Now()

	leaeder_node := oracle.GetTopNLeaders()[0]

	fmt.Printf("attacking leader node %v\n", leaeder_node)

	for j := 0; j < len(nodes); j++ {
		if j == leaeder_node {
			continue
		}
		links[leaeder_node][j].SetDelay(100)
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

	fmt.Print("RaftAttack5 complete\n")
}

type RaftAttack6 struct {
	logger *util.Logger
}

func NewRaftAttack6(logger *util.Logger) *RaftAttack6 {
	return &RaftAttack6{
		logger: logger,
	}
}

func (a *RaftAttack6) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {

	fmt.Printf("Running RaftAttack6 all links have close to view timeout delay for %v seconds\n", duration)

	start_time := time.Now()

	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if i == j {
				continue
			}
			links[i][j].SetDelay(450)
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

	fmt.Print("RaftAttack6 complete\n")
}

type RaftAttack7 struct {
	logger *util.Logger
}

func NewRaftAttack7(logger *util.Logger) *RaftAttack7 {
	return &RaftAttack7{
		logger: logger,
	}
}

func (a *RaftAttack7) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {

	fmt.Printf("Running RaftAttack7 two stragglers for %v seconds\n", duration)

	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-10) {
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		nodes[0].Pause()
		nodes[1].Pause()
		time.Sleep(5 * time.Second)
		nodes[0].Continue()
		nodes[1].Continue()
	}

	fmt.Print("RaftAttack7 complete\n")
}

type RaftAttack8 struct {
	logger *util.Logger
}

func NewRaftAttack8(logger *util.Logger) *RaftAttack8 {
	return &RaftAttack8{
		logger: logger,
	}
}

func (a *RaftAttack8) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {

	fmt.Printf("Running RaftAttack8 straggler in closest majority for %v seconds\n", duration)

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
		//nodes[0].Pause() // straggler in majority
		time.Sleep(1 * time.Second)
		//nodes[0].Continue() // straggler in majority
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

	fmt.Print("RaftAttack8 complete\n")
}
