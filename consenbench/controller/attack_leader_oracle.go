package controller

import (
	"fmt"
	"sort"
	"toture-test/consenbench/common"
	"toture-test/util"
)

type LeaderOracle struct {
	nodes  []*common.Node
	logger *util.Logger
}

func NewLeaderOracle(nodes []*common.Node, logger *util.Logger) *LeaderOracle {
	return &LeaderOracle{
		nodes:  nodes,
		logger: logger,
	}
}

// get the node ids of nodes which consume the most resources

func (l *LeaderOracle) GetTopNLeaders() []int {
	type NodeStats struct {
		node      *common.Node
		cpuSum    float32
		netInSum  float32
		netOutSum float32
	}

	var nodeStatsList []NodeStats

	// Gather stats for each node
	for _, node := range l.nodes {
		cpu_usage, _, network_in, network_out := node.GetStats()

		cpuUsage := cpu_usage[Max(0, len(cpu_usage)-5):]
		netIn := network_in[Max(0, len(network_in)-5):]
		netOut := network_out[Max(0, len(network_out)-5):]

		// Compute the sums of the last 3 slots
		cpuSum := Sum(cpuUsage)
		netInSum := Sum(netIn)
		netOutSum := Sum(netOut)

		nodeStatsList = append(nodeStatsList, NodeStats{
			node:      node,
			cpuSum:    cpuSum / float32(len(cpuUsage)),
			netInSum:  netInSum / float32(len(netIn)),
			netOutSum: netOutSum / float32(len(netOut)),
		})
	}

	// Sort nodes by cpuSum first, then netInSum, and then netOutSum in descending order
	sort.Slice(nodeStatsList, func(i, j int) bool {
		if nodeStatsList[i].cpuSum != nodeStatsList[j].cpuSum {
			return nodeStatsList[i].cpuSum > nodeStatsList[j].cpuSum
		}
		if nodeStatsList[i].netInSum != nodeStatsList[j].netInSum {
			return nodeStatsList[i].netInSum > nodeStatsList[j].netInSum
		}
		return nodeStatsList[i].netOutSum > nodeStatsList[j].netOutSum
	})

	// Calculate how many leaders to select
	numLeaders := len(l.nodes)

	l.logger.Debug(fmt.Sprintf("Node resource usage --"), 3)

	// Extract the IDs of the top n leaders
	var leaderIds []int
	for i := 0; i < numLeaders && i < len(nodeStatsList); i++ {
		leaderIds = append(leaderIds, nodeStatsList[i].node.Id-2) // controller is node 1, and the first replica is not 2
		l.logger.Debug(fmt.Sprintf(fmt.Sprintf("ID: %v, CPU: %v\n", nodeStatsList[i].node.Id-2, nodeStatsList[i].cpuSum)), 0)
	}

	return leaderIds
}
