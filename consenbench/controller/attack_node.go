package controller

import (
	"fmt"
	"toture-test/consenbench/common"
	"toture-test/util"
)

type AttackNode struct {
	Id           int
	Controller   *Controller
	Process_name string
	logger       *util.Logger
}

func NewAttackNode(id int, controller *Controller, name string, logger *util.Logger) *AttackNode {
	return &AttackNode{
		Id:           id,
		Controller:   controller,
		Process_name: name,
		logger:       logger,
	}
}

// send the initial replica specific information to the client

func (n *AttackNode) Init(ports []string, id_ip []string) {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Init),
				StringArgs:    append([]string{n.Process_name}, ports...),
				Ips:           id_ip,
			},
		},
		Peer: n.Id,
	})
	n.logger.Debug(fmt.Sprintf("Init node %v with process name %v and listening ports %v", n.Id, n.Process_name, ports), 3)
}

func (n *AttackNode) Kill() {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Kill),
			},
		},
		Peer: n.Id,
	})
	n.logger.Debug(fmt.Sprintf("killed node %v", n.Id), 3)
}

func (n *AttackNode) Pause() {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Pause),
			},
		},
		Peer: n.Id,
	})
	n.logger.Debug(fmt.Sprintf("Paused node %v", n.Id), 3)
}

func (n *AttackNode) Continue() {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Continue),
			},
		},
		Peer: n.Id,
	})
	n.logger.Debug(fmt.Sprintf("Continued node %v", n.Id), 3)
}

func (n *AttackNode) SetSkew(ms float32) {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().SetSkew),
				FloatArgs:     []float32{ms},
			},
		},
		Peer: n.Id,
	})
	n.logger.Debug(fmt.Sprintf("Skewed node %v", n.Id), 3)
}

func (n *AttackNode) SetDrift(ms float32) {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().SetDrift),
				FloatArgs:     []float32{ms},
			},
		},
		Peer: n.Id,
	})
	n.logger.Debug(fmt.Sprintf("Drfted node %v", n.Id), 3)
}
