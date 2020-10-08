package main

type WaterSectionManager struct {
	Sections  []*WaterSection
	Schedules []*SectionSchedule
}

type waterSectionID int8

func NewWaterSectionManager(sections ...*WaterSection) *WaterSectionManager {
	m := &WaterSectionManager{
		Sections:  make([]*WaterSection, len(sections)),
		Schedules: make([]*SectionSchedule, len(sections)),
	}
	for idx, section := range sections {
		m.Sections[idx] = section
		m.Schedules[idx] = &SectionSchedule{WaterSectionID: waterSectionID(idx)}
	}
	return m
}

func (wsm *WaterSectionManager) Process(now timeUnit) {
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

func (wsm *WaterSectionManager) AllOff(now timeUnit) (anyOn bool) {
	for id, schedule := range wsm.Schedules {
		if schedule.isOn {
			anyOn = true
			wsm.sectionOff(now, waterSectionID(id))
		}
	}
	return
}

func (wsm *WaterSectionManager) ForceOn(idx waterSectionID, now timeUnit) (success bool) {
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

func (wsm *WaterSectionManager) Update(idx waterSectionID, on, off timeUnit, nextActionAt timeUnit) bool {
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		schedule.OnSeconds = on
		schedule.OffSeconds = off
		schedule.NextActionAt = nextActionAt
		return true
	}
	return false
}

func (wsm *WaterSectionManager) UpdateOnSeconds(idx waterSectionID, n timeUnit) bool {
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		schedule.OnSeconds = n
		return true
	}
	return false
}

func (wsm *WaterSectionManager) UpdateOffSeconds(idx waterSectionID, n timeUnit) bool {
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		schedule.OffSeconds = n
		return true
	}
	return false
}

func (wsm *WaterSectionManager) UpdateNextActionAt(idx waterSectionID, n timeUnit) bool {
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

func (wsm *WaterSectionManager) sectionOn(now timeUnit, idx waterSectionID) {
	if int(idx) < len(wsm.Sections) {
		wsm.Sections[idx].On()
		wsm.Schedules[idx].On(now)
	}
}

func (wsm *WaterSectionManager) sectionOff(now timeUnit, idx waterSectionID) {
	if int(idx) < len(wsm.Sections) {
		wsm.Sections[idx].Off()
		wsm.Schedules[idx].Off(now)
	}
}

type SectionSchedule struct {
	OnSeconds, OffSeconds timeUnit
	WaterSectionID        waterSectionID
	// state
	NextActionAt      timeUnit
	isOn              bool
	OnAccum, OffAccum timeUnit
	LastActionAt      timeUnit
}

func (s *SectionSchedule) IsOn() bool {
	return s.isOn
}

func (s *SectionSchedule) ShouldTurnOn(n timeUnit) bool {
	return !s.isOn && s.NextActionAt < n && s.OnSeconds > 0
}

func (s *SectionSchedule) ShouldTurnOff(n timeUnit) bool {
	return s.isOn && s.NextActionAt < n && s.OffSeconds > 0
}

func (s *SectionSchedule) On(n timeUnit) {
	s.NextActionAt = n + timeUnit(s.OnSeconds)
	s.isOn = true
	s.OffAccum += now() - s.LastActionAt
	s.LastActionAt = now()
}

func (s *SectionSchedule) Off(n timeUnit) {
	s.NextActionAt = n + timeUnit(s.OffSeconds)
	s.isOn = false
	s.OnAccum += now() - s.LastActionAt
	s.LastActionAt = now()
}
