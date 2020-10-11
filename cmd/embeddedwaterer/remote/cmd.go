package remote

import (
	"bytes"

	"github.com/compnski/plantwaterer/cmd/embeddedwaterer/schedule"
	"github.com/compnski/plantwaterer/cmd/embeddedwaterer/timing"
)

type CommandStatus uint8

type Command uint8

const (
	StatusOK      CommandStatus = 0
	StatusErr     CommandStatus = 1
	StatusNYI     CommandStatus = 2
	StatusUnknown CommandStatus = 255

	Unused     Command = 1
	CmdTurnOn  Command = 2
	CmdTurnOff Command = 3
	CmdSetOn   Command = 4 // Seconds
	CmdSetOff  Command = 5 // Minutes
	CmdSetNext Command = 6
)

var CmdPrefix = []byte{'C', 'M', 'D'}

type RemoteController struct {
	wsm *schedule.WaterSectionManager /*WaterSystem*/
}

// WaterSystem is sadly unusable because of tinygo type constraints on the AVR.
type WaterSystem interface {
	UpdateOnSeconds(schedule.WaterSectionID, timing.TimeUnit) bool
	UpdateOffSeconds(schedule.WaterSectionID, timing.TimeUnit) bool
	UpdateNextActionAt(schedule.WaterSectionID, timing.TimeUnit) bool
	ForceOn(schedule.WaterSectionID, timing.TimeUnit) bool
	AllOff(now timing.TimeUnit) bool
}

func NewRemoteController(wsm *schedule.WaterSectionManager) *RemoteController {
	return &RemoteController{
		wsm: wsm,
	}
}

func (c *RemoteController) ProcessBytes(now timing.TimeUnit, data []byte) {
	if bytes.HasPrefix(data, CmdPrefix) {
		command, id, value := Command(data[3]), data[4], uint16(data[5])<<8+uint16(data[6])
		status := c.Dispatch(now, command, id, value)
		println(`{"cmd":`, command, `"id":`, id, `"value":`, value, `"status":`, status, `}`)
	} else {
		println(`{"baddata":"`, string(data), `"}`)
	}
}

func (c *RemoteController) Dispatch(now timing.TimeUnit, command Command, id uint8, value uint16) CommandStatus {
	switch command {
	case CmdTurnOn:
		return c.TurnOn(id, now)
	case CmdTurnOff:
		return c.TurnOff(now)
	case CmdSetOn:
		return c.SetOnSeconds(id, value)
	case CmdSetOff:
		return c.SetOffMinutes(id, value)
	case CmdSetNext:
		return c.SetNextActionIn(id, now, value)
	}
	return StatusUnknown
}

func (c *RemoteController) SetOnSeconds(waterSection uint8, n uint16) CommandStatus {
	if c.wsm.UpdateOnSeconds(schedule.WaterSectionID(waterSection), timing.TimeUnit(n)) {
		return StatusOK
	}
	return StatusErr
}

func (c *RemoteController) SetOffMinutes(waterSection uint8, n uint16) CommandStatus {
	if c.wsm.UpdateOffSeconds(schedule.WaterSectionID(waterSection), timing.TimeUnit(n)*60) {
		return StatusOK
	}
	return StatusErr
}

func (c *RemoteController) SetNextActionIn(waterSection uint8, now timing.TimeUnit, nDelta uint16) CommandStatus {
	nextAt := now + timing.TimeUnit(nDelta)
	if c.wsm.UpdateNextActionAt(schedule.WaterSectionID(waterSection), nextAt) {
		return StatusOK
	}
	return StatusErr
}

func (c *RemoteController) TurnOn(waterSection uint8, now timing.TimeUnit) CommandStatus {
	if c.wsm.ForceOn(schedule.WaterSectionID(waterSection), now) {
		return StatusOK
	}
	return StatusErr
}

func (c *RemoteController) TurnOff(now timing.TimeUnit) CommandStatus {
	if c.wsm.AllOff(now) {
		return StatusOK
	}
	return StatusErr

}
