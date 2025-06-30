package controller

//
//import (
//	"fmt"
//	"math/rand"
//	"time"
//	"toture-test/util"
//)
//
//// DriftAttack_1 is an attack that changes the drift of the leader node to all other nodes
//
//type DriftAttack_1 struct {
//	logger *util.Logger
//}
//
//func NewDriftAttack_1(logger *util.Logger) *DriftAttack_1 {
//	return &DriftAttack_1{
//		logger: logger,
//	}
//}
//
//func (a *DriftAttack_1) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
//	fmt.Printf("Running Drift Attack_1 : changing the drift of leader for %v seconds\n", duration)
//	start_time := time.Now()
//
//	drifetedNodes := make(map[int]bool)
//	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
//		node_id := oracle.GetTopNLeaders()[0]
//		// set all outgoing links to of leader to high delay
//		if !drifetedNodes[node_id] {
//			drifetedNodes[node_id] = true
//		} else {
//			continue
//		}
//		fmt.Printf("attacking leader node %v\n", node_id)
//		for i := 0; i < len(nodes); i++ {
//			for j := 0; j < len(nodes); j++ {
//				if i == j {
//					continue
//				}
//				if i == node_id {
//					go drift(links[i][j], float64(duration-15)-time.Now().Sub(start_time).Seconds())
//					fmt.Printf("setting drifting %v\n", i)
//				}
//			}
//		}
//
//		time.Sleep(3 * time.Second)
//	}
//
//	fmt.Print(" attack complete\n")
//}
//
//// DriftAttack_2 is an attack that changes the drift of the minority node to all other nodes
//
//type DriftAttack_2 struct {
//	logger *util.Logger
//}
//
//func NewDriftAttack_2(logger *util.Logger) *DriftAttack_2 {
//	return &DriftAttack_2{
//		logger: logger,
//	}
//}
//
//func (a *DriftAttack_2) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
//	fmt.Printf("Running Drift Attack_2 : changing the drift of minority for %v seconds\n", duration)
//	start_time := time.Now()
//	drifetedNodes := make(map[int]bool)
//	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
//		randNodes := make(map[int]bool)
//
//		majority := len(nodes)/2 + 1
//		minority := len(nodes) - majority
//
//		for i := 0; i < minority; i++ {
//			rand_node := rand.Intn(len(nodes))
//			for randNodes[rand_node] {
//				rand_node = rand.Intn(len(nodes))
//			}
//			randNodes[rand_node] = true
//		}
//		fmt.Printf("attacking minority nodes %v \n", randNodes)
//		for i := 0; i < len(nodes); i++ {
//			if drifetedNodes[i] {
//				continue
//			}
//			for j := 0; j < len(nodes); j++ {
//				if i == j {
//					continue
//				}
//
//				if randNodes[i] {
//					go drift(links[i][j], float64(duration-15)-time.Now().Sub(start_time).Seconds())
//					fmt.Printf("setting drifting %v\n", i)
//				}
//
//			}
//			drifetedNodes[i] = true
//		}
//
//		time.Sleep(3 * time.Second)
//	}
//
//	fmt.Print(" attack complete\n")
//}
//
//// DriftAttack_3 is an attack that changes the drift of the majority node to all other nodes
//
//type DriftAttack_3 struct {
//	logger *util.Logger
//}
//
//func NewDriftAttack_3(logger *util.Logger) *DriftAttack_3 {
//	return &DriftAttack_3{
//		logger: logger,
//	}
//}
//
//func (a *DriftAttack_3) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
//	fmt.Printf("Running Drift Attack_3 : changing the drift of majorty for %v seconds\n", duration)
//	start_time := time.Now()
//	drifetedNodes := make(map[int]bool)
//	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {
//		randNodes := make(map[int]bool)
//
//		majority := len(nodes)/2 + 1
//
//		for i := 0; i < majority; i++ {
//			rand_node := rand.Intn(len(nodes))
//			for randNodes[rand_node] {
//				rand_node = rand.Intn(len(nodes))
//			}
//			randNodes[rand_node] = true
//		}
//		fmt.Printf("attacking majority nodes %v \n", randNodes)
//		for i := 0; i < len(nodes); i++ {
//			if drifetedNodes[i] {
//				continue
//			}
//			for j := 0; j < len(nodes); j++ {
//				if i == j {
//					continue
//				}
//				if randNodes[i] {
//					go drift(links[i][j], float64(duration-15)-time.Now().Sub(start_time).Seconds())
//					fmt.Printf("setting drifting %v\n", i)
//				}
//			}
//			drifetedNodes[i] = true
//		}
//
//		time.Sleep(3 * time.Second)
//	}
//
//	fmt.Print("attack complete\n")
//}
//
//func drift(link *AttackLink, remaining_seconds float64) {
//	start_time := time.Now()
//	drift := 10
//	for time.Now().Sub(start_time).Seconds() < remaining_seconds {
//		link.SetDelay(float32(drift))
//		drift = drift + 10
//	}
//}
