package main

import "io"

type CommandStatus uint8

type Command_Command uint8

const (
	StatusOK      CommandStatus = 0
	StatusErr     CommandStatus = 1
	StatusNYI     CommandStatus = 2
	StatusUnknown CommandStatus = 255

	CmdStats Command_Command = iota
	CmdTurnOn
	CmdTurnOff
	CmdSetOn
	CmdSetOff
	CmdSetNext
	CmdZeroStats

	CmdMask   = 0xF000
	IdMask    = 0x0F00
	ValueMask = 0x00FF
)

type Controller struct {
	wsm *WaterSectionManager
}

func (c *Controller) ProcessCommand(commandBits uint32) CommandStatus {
	var command = Command_Command((CmdMask & commandBits) >> 24)
	var id = uint8((IdMask & commandBits) >> 16)
	var value = uint16((ValueMask & commandBits))
	var writer io.Writer

	switch command {
	case CmdStats:
		return c.SendStats(writer)
	case CmdTurnOn:
		return c.SetOnSeconds(id, value)
	case CmdTurnOff:
		return c.SetOffSeconds(id, value)
	case CmdSetOn:
		return c.TurnOn(id)
	case CmdSetOff:
		return c.TurnOff()
	case CmdSetNext:
		return c.SetNextActionIn(id, value)
	case CmdZeroStats:
		return c.ZeroStats()
	}
	return StatusUnknown
}

func (c *Controller) SetOnSeconds(waterSection uint8, n uint16) CommandStatus {
	if c.wsm.UpdateOnSeconds(waterSectionID(waterSection), n) {
		return StatusOK
	}
	return StatusErr
}

func (c *Controller) SetOffSeconds(waterSection uint8, n uint16) CommandStatus {
	if c.wsm.UpdateOffSeconds(waterSectionID(waterSection), n) {
		return StatusOK
	}
	return StatusErr
}

func (c *Controller) SetNextActionIn(waterSection uint8, nDelta uint16) CommandStatus {
	// TODO: Needs now?
	nextAt := seconds + timeUnit(nDelta)
	if c.wsm.UpdateNextActionAt(waterSectionID(waterSection), nextAt) {
		return StatusOK
	}
	return StatusErr
}

func (c *Controller) TurnOn(waterSection uint8) CommandStatus {
	// TODO: Needs now
	return StatusNYI
	// if c.wsm.ForceOn(waterSectionID(waterSection)) {
	// 		return StatusOK
	// 	}
	// 	return StatusErr
}

func (c *Controller) TurnOff() CommandStatus {
	// TODO: Needs now
	return StatusNYI
	// if c.wsm.AllOff(waterSectionID(waterSection)) {
	// 	return StatusOK
	// }
	// return StatusErr

}

func (c *Controller) SendStats(w io.Writer) CommandStatus {
	return StatusNYI
}

func (c *Controller) ZeroStats() CommandStatus {
	return StatusNYI
}
