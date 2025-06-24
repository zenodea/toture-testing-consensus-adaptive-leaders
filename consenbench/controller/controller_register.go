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
	case "skew_1":
		return NewSkewAttack_1(c.logger)
	case "skew_2":
		return NewSkewAttack_2(c.logger)
	case "skew_3":
		return NewSkewAttack_3(c.logger)
	case "drift_1":
		return NewDriftAttack_1(c.logger)
	case "drift_2":
		return NewDriftAttack_2(c.logger)
	case "drift_3":
		return NewDriftAttack_3(c.logger)
	case "leader_1":
		return NewLeaderAttack1(c.logger)
	case "leader_2":
		return NewLeaderAttack2(c.logger)
	case "leader_3":
		return NewLeaderAttack3(c.logger)
	case "leader_4":
		return NewLeaderAttack4(c.logger)
	case "leader_5":
		return NewLeaderAttack5(c.logger)
	case "leader_6":
		return NewLeaderAttack6(c.logger)
	case "leader_7":
		return NewLeaderAttack7(c.logger)
	case "leader_8":
		return NewLeaderAttack8(c.logger)
	case "leader_9":
		return NewLeaderAttack9(c.logger)
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
