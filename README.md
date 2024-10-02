# Torture-testing-consensus

This repository implements a networked tool to test the robustness of consensus algorithms.
It allows plugging in different consensus algorithms and network topologies to test their performance under different conditions.

## NOTES

This project is currently under **heavy development** and is **not yet ready for production use**.

## Indexing

*consenbench/assets/ip.yaml* contains 1 indexed nodes, and node with id 1 will be the controller and replica and client nodes will be indexed starting from 2

In the attack interface *Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int)* :

- nodes: list of replica nodes in the network, with 0 indexing; nodes[0] is the first replica with a consensus replica running in it (id = 2)
- links: list of links between nodes, with 0 indexing; links[0][1] is the link between node id 2 and 3
- oracle: LeaderOracle retrieves the set of leaders sorted in the descending order of their resource usage, and contains node index of the nodes array