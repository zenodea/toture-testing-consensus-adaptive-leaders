# Torture-testing-consensus

This repository implements a networked tool to test the robustness of consensus algorithms.
It allows plugging in different consensus algorithms and network topologies to test their performance under different scenarios.

## Supported attacks

- Bandwidth throttling
- Crashes
- Delay injection
- Network Partitions
- Clock skew
- Clock drift
- Straggler nodes

## Supported consensus algorithms

- Raft
- Paxos (2 versions)
- Baxos
- Rabia
- RACS
- SADL-RACS
- EPaxos
- Mencius
- Generalized Paxos
- QuePaxa
- Mahi-Mahi
- Codial-Miners
- Mysticeti
- Jolteon
- HotStuff