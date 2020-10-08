package main

import (
	"bytes"
	"io"
)

type CommandStatus uint8

type Command uint8

const (
	StatusOK      CommandStatus = 0
	StatusErr     CommandStatus = 1
	StatusNYI     CommandStatus = 2
	StatusUnknown CommandStatus = 255

	CmdStats     Command = 1
	CmdTurnOn    Command = 2
	CmdTurnOff   Command = 3
	CmdSetOn     Command = 4 // Seconds
	CmdSetOff    Command = 5 // Minutes
	CmdSetNext   Command = 6
	CmdZeroStats Command = 7

	CmdMask   = 0xF000
	IdMask    = 0x0F00
	ValueMask = 0x00FF
)

type Controller struct {
	wsm *WaterSectionManager
}

func (c *Controller) ProcessBytes(data []byte) {
	//println(string(data), data[0], data[1], data[2], data[3], data[4], data[5], data[6])
	if bytes.HasPrefix(data, CmdPrefix) {
		//i := uint32(data[3])<<24 + uint32(data[4])<<16 + uint32(data[5])<<8 + uint32(data[6])
		//println(i)
		command, id, value := Command(data[3]), data[4], uint16(data[5])<<8+uint16(data[6])
		status := c.Dispatch(command, id, value)
		println(`{"cmd":`, command, `"id":`, id, `"value":`, value, `"status":`, status, `}`)
	} else {
		println(`{"baddata":"`, string(data), `"}`)
	}
}

//func (c *Controller) ProcessCommand(commandBits uint32) (command Command,
//	id uint8,
//value uint16,
//status CommandStatus) {

func (c *Controller) Dispatch(command Command, id uint8, value uint16) CommandStatus {
	switch command {
	case CmdTurnOn:
		//		println("on", id)
		return c.TurnOn(id)
	case CmdTurnOff:
		//println("off", id)
		return c.TurnOff()
	case CmdSetOn:
		//println("son", id, value)
		return c.SetOnSeconds(id, value)
	case CmdSetOff:
		//println("sof", id, value)
		return c.SetOffMinutes(id, value)
	case CmdSetNext:
		//println("sn", id, value)
		return c.SetNextActionIn(id, value)
	case CmdZeroStats:
		return c.ZeroStats()
	}
	return StatusUnknown
}

func (c *Controller) SetOnSeconds(waterSection uint8, n uint16) CommandStatus {
	if c.wsm.UpdateOnSeconds(waterSectionID(waterSection), timeUnit(n)) {
		return StatusOK
	}
	return StatusErr
}

func (c *Controller) SetOffMinutes(waterSection uint8, n uint16) CommandStatus {
	if c.wsm.UpdateOffSeconds(waterSectionID(waterSection), timeUnit(n)*60) {
		return StatusOK
	}
	return StatusErr
}

func (c *Controller) SetNextActionIn(waterSection uint8, nDelta uint16) CommandStatus {
	nextAt := now() + timeUnit(nDelta)
	if c.wsm.UpdateNextActionAt(waterSectionID(waterSection), nextAt) {
		return StatusOK
	}
	return StatusErr
}

func (c *Controller) TurnOn(waterSection uint8) CommandStatus {
	if c.wsm.ForceOn(waterSectionID(waterSection), now()) {
		return StatusOK
	}
	return StatusErr
}

func (c *Controller) TurnOff() CommandStatus {
	if c.wsm.AllOff(now()) {
		return StatusOK
	}
	return StatusErr

}

func (c *Controller) SendStats(w io.Writer) CommandStatus {
	return StatusNYI
}

func (c *Controller) ZeroStats() CommandStatus {
	return StatusNYI
}
