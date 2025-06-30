package controller

import (
	"toture-test/protocols"
	bullshark "toture-test/protocols/bullshark"
	cft "toture-test/protocols/cft_dag"
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
	case "noop":
		return NewNoopAttack(c.logger)
	case "LeaderPartition":
		return NewLeaderPartition(c.logger)
	case "OneQuorumNodePartition":
		return NewOneQuorumNodePartition(c.logger)
	case "MajorityHighDelay":
		return NewMajorityHighDelay(c.logger)
	case "MinorityCrash":
		return NewMinorityCrash(c.logger)
	case "MinorityStraggler":
		return NewMinorityStraggler(c.logger)
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
	} else if protocol == "cft_dag" {
		return cft.NewCFT_DAG(c.logger)
	} else {
		panic("Unknown protocol")
	}
}
