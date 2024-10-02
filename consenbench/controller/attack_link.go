package controller

import (
	"fmt"
	"toture-test/consenbench/common"
	"toture-test/util"
)

type AttackLink struct {
	Id_sender   int
	Id_reciever int
	Controller  *Controller
	logger      *util.Logger
}

func NewAttackLink(sender int, reciever int, controller *Controller, logger *util.Logger) *AttackLink {
	return &AttackLink{
		Id_sender:   sender,
		Id_reciever: reciever,
		Controller:  controller,
		logger:      logger,
	}
}

func (a *AttackLink) SetStatus(on bool) {
	if !on {
		a.SetLoss(100)
	} else {
		a.SetLoss(0)
	}
	a.logger.Debug(fmt.Sprintf("set status link %v-%v", a.Id_sender, a.Id_reciever), 3)
}

func (a *AttackLink) SetDelay(ms float32) {
	a.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().SetDelay),
				FloatArgs:     []float32{ms},
				IntArgs:       []int32{int32(a.Id_reciever)},
			},
		},
		Peer: a.Id_sender,
	})
	a.logger.Debug(fmt.Sprintf("set delay link %v-%v", a.Id_sender, a.Id_reciever), 3)
}

func (a *AttackLink) SetLoss(l float32) {
	if l == 100 {
		l = 95 // to avoid 100% loss which will make the testbed unresponsive
	}
	a.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().SetLoss),
				FloatArgs:     []float32{l},
				IntArgs:       []int32{int32(a.Id_reciever)},
			},
		},
		Peer: a.Id_sender,
	})
	a.logger.Debug(fmt.Sprintf("set loss link %v-%v", a.Id_sender, a.Id_reciever), 3)
}

func (a *AttackLink) SetBandwidth(b float32) {
	a.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().SetBandwidth),
				FloatArgs:     []float32{b},
				IntArgs:       []int32{int32(a.Id_reciever)},
			},
		},
		Peer: a.Id_sender,
	})
	a.logger.Debug(fmt.Sprintf("set bandwidth link %v-%v", a.Id_sender, a.Id_reciever), 3)

}
