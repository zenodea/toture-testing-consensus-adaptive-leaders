package client

// pause the client

func (c *Client) Pause() {
	RunCommand("pkill", []string{"-STOP", "-f", c.Attacker.Process_name}, c.logger)
	c.logger.Debug("paused", 3)
}

// continue the client

func (c *Client) Continue() {
	RunCommand("pkill", []string{"-CONT", "-f", c.Attacker.Process_name}, c.logger)
	c.logger.Debug("continue", 3)
}

// kill the client

func (c *Client) Kill() {
	for _, v := range c.Attacker.NetEmAttackers {
		v.ExecuteLastNetEmCommands()
	}
	c.CleanUp()
	RunCommand("pkill", []string{"-KILL", "-f", c.Attacker.Process_name}, c.logger)
	c.logger.Debug("killed consensus node", 3)
}

// set the skew

func (c *Client) SetSkew(f float32) {
	// TODO
	panic("Not implemented")
}

// set the drift

func (c *Client) SetDrift(f float32) {
	panic("Not implemented")
}
