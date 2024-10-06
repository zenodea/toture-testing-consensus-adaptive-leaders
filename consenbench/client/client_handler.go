package client

import (
	"os"
	"toture-test/consenbench/common"
)

func (c *Client) Handle(msg *common.ControlMsg) {
	if int(msg.OperationType) == common.GetOperationCodes().ShutDown {
		c.logger.Debug("Received ShutDown signal from controller", 3)
		c.CleanUp()
		os.Exit(0)
	} else if int(msg.OperationType) == common.GetOperationCodes().Stats {
		panic("Received Stats signal from controller")
	} else if int(msg.OperationType) == common.GetOperationCodes().Kill {
		c.logger.Debug("Received Kill signal from controller", 3)
		c.Kill()
	} else if int(msg.OperationType) == common.GetOperationCodes().Pause {
		c.logger.Debug("Received Pause signal from controller", 3)
		c.Pause()
	} else if int(msg.OperationType) == common.GetOperationCodes().Continue {
		c.logger.Debug("Received Continue signal from controller", 3)
		c.Continue()
	} else if int(msg.OperationType) == common.GetOperationCodes().SetSkew {
		c.logger.Debug("Received SetSkew signal from controller", 3)
		c.SetSkew(msg.FloatArgs[0])
	} else if int(msg.OperationType) == common.GetOperationCodes().SetDrift {
		c.logger.Debug("Received SetDrift signal from controller", 3)
		c.SetDrift(msg.FloatArgs[0])
	} else if int(msg.OperationType) == common.GetOperationCodes().SetDelay {
		c.logger.Debug("Received SetDelay signal from controller", 3)
		c.SetDelay(msg.FloatArgs[0], msg.IntArgs[0])
	} else if int(msg.OperationType) == common.GetOperationCodes().SetLoss {
		c.logger.Debug("Received SetLoss signal from controller", 0)
		c.SetLoss(msg.FloatArgs[0], msg.IntArgs[0])
	} else if int(msg.OperationType) == common.GetOperationCodes().SetBandwidth {
		c.logger.Debug("Received SetBandwidth signal from controller", 3)
		c.SetBandwidth(msg.FloatArgs[0], msg.IntArgs[0])
	} else if int(msg.OperationType) == common.GetOperationCodes().Init {
		c.logger.Debug("Received Init signal from controller", 3)
		c.InitAttacker(msg)
	} else {
		panic("Unknown operation type")
	}
	// todo handle other messages
}
