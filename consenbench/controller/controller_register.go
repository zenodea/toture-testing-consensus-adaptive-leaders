package controller

import (
	"toture-test/protocols"
	bullshark "toture-test/protocols/bullshark"
	dedis_paxos "toture-test/protocols/dedis_paxos"
	dedis_raft "toture-test/protocols/dedis_raft"
	efficient "toture-test/protocols/efficient"
	etcd "toture-test/protocols/etcd"
	hotstuff_2 "toture-test/protocols/hotstuff_2"
	hotstuff_3 "toture-test/protocols/hotstuff_3"
	mahi "toture-test/protocols/mahi"
	mysticeti "toture-test/protocols/mysticeti"
	ping "toture-test/protocols/ping"
	quepaxa "toture-test/protocols/quepaxa"
	rabia "toture-test/protocols/rabia"
	racs "toture-test/protocols/racs"
	sadl_racs "toture-test/protocols/sadl_racs"
	tusk "toture-test/protocols/tusk"
	zoo_keeper "toture-test/protocols/zoo_keeper"
)

func (c *Controller) GetAttackImpl() Attack {
	switch c.Options.Attack {
	case "basic":
		return NewBasicAttack(c.logger)
	case "noop":
		return NewNoopAttack(c.logger)
	case "partition_1":
		return NewPartitionAttack_1(c.logger)
	case "partition_2":
		return NewPartitionAttack_2(c.logger)
	case "partition_3":
		return NewPartitionAttack_3(c.logger)
	case "partition_4":
		return NewPartitionAttack_4(c.logger)
	case "partition_5":
		return NewPartitionAttack_5(c.logger)
	case "partition_6":
		return NewPartitionAttack_6(c.logger)
	case "partition_7":
		return NewPartitionAttack_7(c.logger)
	case "delay_1":
		return NewDelayAttack_1(c.logger)
	case "delay_2":
		return NewDelayAttack_2(c.logger)
	case "delay_3":
		return NewDelayAttack_3(c.logger)
	case "delay_4":
		return NewDelayAttack_4(c.logger)
	case "bandwidth_1":
		return NewBandwidthAttack_1(c.logger)
	case "bandwidth_2":
		return NewBandwidthAttack_2(c.logger)
	case "bandwidth_3":
		return NewBandwidthAttack_3(c.logger)
	case "bandwidth_4":
		return NewBandwidthAttack_4(c.logger)
	case "skew_1":
		return NewSkewAttack_1(c.logger)
	case "skew_2":
		return NewSkewAttack_2(c.logger)
	case "skew_3":
		return NewSkewAttack_3(c.logger)
	case "straggler_1":
		return NewStragglerAttack_1(c.logger)
	case "straggler_2":
		return NewStragglerAttack_2(c.logger)
	case "straggler_3":
		return NewStragglerAttack_3(c.logger)
	case "crash_1":
		return NewCrashAttack_1(c.logger)
	case "crash_2":
		return NewCrashAttack_2(c.logger)
	case "crash_3":
		return NewCrashAttack_3(c.logger)
	case "drift_1":
		return NewDriftAttack_1(c.logger)
	case "drift_2":
		return NewDriftAttack_2(c.logger)
	case "drift_3":
		return NewDriftAttack_3(c.logger)
	case "raft_1":
		return NewRaftAttack1(c.logger)
	case "raft_2":
		return NewRaftAttack2(c.logger)
	case "raft_3":
		return NewRaftAttack3(c.logger)
	case "raft_4":
		return NewRaftAttack4(c.logger)
	case "raft_5":
		return NewRaftAttack5(c.logger)
	case "raft_6":
		return NewRaftAttack6(c.logger)
	case "raft_7":
		return NewRaftAttack7(c.logger)
	case "raft_8":
		return NewRaftAttack8(c.logger)
	default:
		panic("Unknown attack: " + c.Options.Attack)
	}
}

// as you add more protocols, you need to add the protocol here

func (c *Controller) GetProtocolImpl(protocol string) protocols.Consensus {
	if protocol == "ping" {
		return ping.NewPing(c.logger)
	} else if protocol == "dedis_paxos" {
		return dedis_paxos.NewDedis_Paxos(c.logger)
	} else if protocol == "dedis_raft" {
		return dedis_raft.NewDedis_Raft(c.logger)
	} else if protocol == "rabia" {
		return rabia.NewRabia(c.logger)
	} else if protocol == "sadl_racs" {
		return sadl_racs.NewSadl_Racs(c.logger)
	} else if protocol == "racs" {
		return racs.NewRacs(c.logger)
	} else if protocol == "efficient" {
		return efficient.NewEfficient(c.logger)
	} else if protocol == "quepaxa" {
		return quepaxa.NewQuePaxa(c.logger)
	} else if protocol == "mahi" {
		return mahi.NewMahi(c.logger)
	} else if protocol == "mysticeti" {
		return mysticeti.NewMysticeti(c.logger)
	} else if protocol == "hotstuff_2" {
		return hotstuff_2.NewHotstuff_2(c.logger)
	} else if protocol == "hotstuff_3" {
		return hotstuff_3.NewHotstuff_3(c.logger)
	} else if protocol == "tusk" {
		return tusk.NewTusk(c.logger)
	} else if protocol == "bullshark" {
		return bullshark.NewBullshark(c.logger)
	} else if protocol == "etcd" {
		return etcd.NewETCD(c.logger)
	} else if protocol == "zoo_keeper" {
		return zoo_keeper.NewZooKeeper(c.logger)
	} else {
		panic("Unknown protocol")
	}
}
