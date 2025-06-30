package controller

import (
	"fmt"
	"time"
	"toture-test/util"
)

type NoopAttack struct {
	logger *util.Logger
}

func NewNoopAttack(logger *util.Logger) *NoopAttack {
	return &NoopAttack{
		logger: logger,
	}
}

func (a *NoopAttack) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-5) {
		a.logger.Debug(fmt.Sprintf("The leader order is %v\n", oracle.GetTopNLeaders()), 0)
		time.Sleep(1 * time.Second)
	}

	a.logger.Debug(fmt.Sprintf("Noop attack complete\n"), 0)
}
