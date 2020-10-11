package schedule

import (
	"github.com/compnski/plantwaterer/cmd/embeddedwaterer/hardware"
	"github.com/compnski/plantwaterer/cmd/embeddedwaterer/timing"
	"github.com/compnski/plantwaterer/cmd/embeddedwaterer/util"
)

type WaterSectionManager struct {
	Sections  []*hardware.WaterSection
	Schedules []*SectionSchedule
}

type WaterSectionID int8

func NewWaterSectionManager(sections ...*hardware.WaterSection) *WaterSectionManager {
	m := &WaterSectionManager{
		Sections:  make([]*hardware.WaterSection, len(sections)),
		Schedules: make([]*SectionSchedule, len(sections)),
	}
	for idx, section := range sections {
		m.Sections[idx] = section
		m.Schedules[idx] = &SectionSchedule{WaterSectionID: WaterSectionID(idx)}
	}
	return m
}

func (wsm *WaterSectionManager) SendStats(now timing.TimeUnit) {
	println("[")
	for idx := range wsm.Sections {
		var (
			schedule = wsm.Schedules[idx]
			accDelta = now - schedule.LastActionAt
			onAccum  = schedule.OnAccum
			offAccum = schedule.OffAccum
		)
		// [on|off]Accum only changes on state change, add the time since last
		// action to smooth constant progress rather than large jumps.
		if schedule.IsOn() {
			onAccum += accDelta
		} else {
			offAccum += accDelta
		}
		println(util.MaybeLeadingComma(idx), `{"id":`, idx, ",",
			`"on":`, schedule.IsOn(), ",",
			`"next":`, schedule.NextActionAt, ",",
			`"last":`, schedule.LastActionAt, ",",
			`"onTime":`, schedule.OnSeconds, ",",
			`"offTime":`, schedule.OffSeconds, ",",
			`"onAcc":`, onAccum, ",",
			`"offAcc":`, offAccum,
			`}`,
		)
	}
	println("]")
}

func (wsm *WaterSectionManager) Process(now timing.TimeUnit) {
	if schedule := wsm.IsOn(); schedule != nil {
		if schedule.ShouldTurnOff(now) {
			wsm.sectionOff(now, schedule.WaterSectionID)
		}
	} else if schedule := wsm.NextChange(); schedule != nil {
		if schedule.ShouldTurnOn(now) {
			wsm.sectionOn(now, schedule.WaterSectionID)
		}
	}
}

func (wsm *WaterSectionManager) AllOff(now timing.TimeUnit) (anyOn bool) {
	for id, schedule := range wsm.Schedules {
		if schedule.isOn {
			anyOn = true
			wsm.sectionOff(now, WaterSectionID(id))
		}
	}
	return
}

func (wsm *WaterSectionManager) ForceOn(idx WaterSectionID, now timing.TimeUnit) (success bool) {
	if wsm.IsOn() != nil {
		return
	}
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		if !schedule.isOn {
			wsm.sectionOn(now, idx)
			success = true
		}
	}
	return
}

func (wsm *WaterSectionManager) Update(idx WaterSectionID, on, off timing.TimeUnit, nextActionAt timing.TimeUnit) bool {
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		schedule.OnSeconds = on
		schedule.OffSeconds = off
		schedule.NextActionAt = nextActionAt
		return true
	}
	return false
}

func (wsm *WaterSectionManager) UpdateOnSeconds(idx WaterSectionID, n timing.TimeUnit) bool {
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		schedule.OnSeconds = n
		return true
	}
	return false
}

func (wsm *WaterSectionManager) UpdateOffSeconds(idx WaterSectionID, n timing.TimeUnit) bool {
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		schedule.OffSeconds = n
		return true
	}
	return false
}

func (wsm *WaterSectionManager) UpdateNextActionAt(idx WaterSectionID, n timing.TimeUnit) bool {
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		schedule.NextActionAt = n
		return true
	}
	return false
}

func (wsm *WaterSectionManager) NextChange() *SectionSchedule {
	if len(wsm.Schedules) < 1 {
		return nil
	}
	nextSchedule := wsm.Schedules[0]
	for _, s := range wsm.Schedules {
		if s.NextActionAt < nextSchedule.NextActionAt {
			nextSchedule = s
		}
	}
	return nextSchedule
}

func (wsm *WaterSectionManager) IsOn() *SectionSchedule {
	for _, s := range wsm.Schedules {
		if s.isOn {
			return s
		}
	}
	return nil
}

func (wsm *WaterSectionManager) sectionOn(now timing.TimeUnit, idx WaterSectionID) {
	if int(idx) < len(wsm.Sections) {
		wsm.Sections[idx].On()
		wsm.Schedules[idx].On(now)
	}
}

func (wsm *WaterSectionManager) sectionOff(now timing.TimeUnit, idx WaterSectionID) {
	if int(idx) < len(wsm.Sections) {
		wsm.Sections[idx].Off()
		wsm.Schedules[idx].Off(now)
	}
}

type SectionSchedule struct {
	OnSeconds, OffSeconds timing.TimeUnit
	WaterSectionID        WaterSectionID
	// state
	NextActionAt      timing.TimeUnit
	isOn              bool
	OnAccum, OffAccum timing.TimeUnit
	LastActionAt      timing.TimeUnit
}

func (s *SectionSchedule) IsOn() bool {
	return s.isOn
}

func (s *SectionSchedule) ShouldTurnOn(now timing.TimeUnit) bool {
	return !s.isOn && s.NextActionAt < now && s.OnSeconds > 0
}

func (s *SectionSchedule) ShouldTurnOff(now timing.TimeUnit) bool {
	return s.isOn && s.NextActionAt < now && s.OffSeconds > 0
}

func (s *SectionSchedule) On(now timing.TimeUnit) {
	s.NextActionAt = now + timing.TimeUnit(s.OnSeconds)
	s.isOn = true
	s.OffAccum += now - s.LastActionAt
	s.LastActionAt = now
}

func (s *SectionSchedule) Off(now timing.TimeUnit) {
	s.NextActionAt = now + timing.TimeUnit(s.OffSeconds)
	s.isOn = false
	s.OnAccum += now - s.LastActionAt
	s.LastActionAt = now
}
