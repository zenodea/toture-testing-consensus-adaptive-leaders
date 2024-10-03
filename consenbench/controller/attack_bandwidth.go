package controller

import (
	"fmt"
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
	fmt.Printf("Running  bandwidth attack 1  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {

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
	fmt.Printf("Running  bandwidth attack 2  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {

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
	fmt.Printf("Running  bandwidth attack 3  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {

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
	fmt.Printf("Running  bandwidth attack 4  for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-15) {

	}

	fmt.Print(" attack complete\n")
}
