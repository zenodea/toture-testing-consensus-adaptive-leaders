package controller

import (
	"fmt"
	"math/rand"
	"time"
	"toture-test/util"
)

// leader is a Crash

type CrashAttack_1 struct {
	logger *util.Logger
}

func NewCrashAttack_1(logger *util.Logger) *CrashAttack_1 {
	return &CrashAttack_1{
		logger: logger,
	}
}

func (a *CrashAttack_1) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running Crash attack 1: Crash leader for %v seconds\n", duration)

	time.Sleep(10 * time.Second)
	leader := oracle.GetTopNLeaders()[0]
	fmt.Printf("attacking leader node %v\n", leader)
	nodes[leader].Kill()
	fmt.Print(" attack complete\n")
}

// leader 1 and 2 are Crashs

type CrashAttack_2 struct {
	logger *util.Logger
}

func NewCrashAttack_2(logger *util.Logger) *CrashAttack_2 {
	return &CrashAttack_2{
		logger: logger,
	}
}

func (a *CrashAttack_2) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running Crash attack 2: Crash leader 1 and 2 for %v seconds\n", duration)
	time.Sleep(10 * time.Second)
	leader1 := oracle.GetTopNLeaders()[0]
	leader2 := oracle.GetTopNLeaders()[1]
	fmt.Printf("attacking leader nodes %v %v\n", leader1, leader2)
	nodes[leader1].Kill()
	time.Sleep(5 * time.Second)
	nodes[leader2].Kill()

	fmt.Print(" attack complete\n")
}

// minority Crashs

type CrashAttack_3 struct {
	logger *util.Logger
}

func NewCrashAttack_3(logger *util.Logger) *CrashAttack_3 {
	return &CrashAttack_3{
		logger: logger,
	}
}

func (a *CrashAttack_3) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running Crash attack 3: minority Crashs for %v seconds\n", duration)

	time.Sleep(10 * time.Second)

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
		nodes[node_id].Kill()
		time.Sleep(5 * time.Second)
	}

	fmt.Print(" attack complete\n")
}
