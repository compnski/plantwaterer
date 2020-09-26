package main

type SectionSchedule struct {
	OnSeconds, OffSeconds timeUnit
	WaterSectionID        waterSectionID
	// state
	NextActionAt timeUnit
	isOn         bool
}

func (s *SectionSchedule) IsOn() bool {
	return s.isOn
}

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

func (s *SectionSchedule) ShouldTurnOn(n timeUnit) bool {
	return !s.isOn && s.NextActionAt < n && s.OnSeconds > 0
}

func (s *SectionSchedule) ShouldTurnOff(n timeUnit) bool {
	return s.isOn && s.NextActionAt < n && s.OffSeconds > 0
}

func (s *SectionSchedule) On(n timeUnit) {
	s.NextActionAt = n + s.OnSeconds
	s.isOn = true
}

func (s *SectionSchedule) Off(n timeUnit) {
	s.NextActionAt = n + s.OffSeconds
	s.isOn = false
}

func (wsm *WaterSectionManager) SectionOn(idx waterSectionID) {
	if int(idx) < len(wsm.Sections) {
		wsm.Sections[idx].On()
	}
}

func (wsm *WaterSectionManager) SectionOff(idx waterSectionID) {
	if int(idx) < len(wsm.Sections) {
		wsm.Sections[idx].Off()
	}
}

func (wsm *WaterSectionManager) Process(n timeUnit) {
	if schedule := wsm.IsOn(); schedule != nil {
		if schedule.ShouldTurnOff(n) {
			schedule.Off(n)
			wsm.SectionOff(schedule.WaterSectionID)
		}
	} else if schedule := wsm.NextChange(); schedule != nil {
		if schedule.ShouldTurnOn(n) {
			schedule.On(n)
			wsm.SectionOn(schedule.WaterSectionID)
		}
	}
}

func (wsm *WaterSectionManager) AllOff(now timeUnit) (anyOn bool) {
	for id, schedule := range wsm.Schedules {
		if schedule.isOn {
			anyOn = true
			schedule.Off(now)
			wsm.SectionOff(waterSectionID(id))
		}
	}
	return
}

func (wsm *WaterSectionManager) ForceOn(idx waterSectionID, now timeUnit) (success bool) {
	if wsm.IsOn != nil {
		return
	}
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		if !schedule.isOn {
			schedule.On(now)
			wsm.SectionOn(idx)
		}
	}
	return
}

func (wsm *WaterSectionManager) Update(idx waterSectionID, on, off, nextActionAt timeUnit) {
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		schedule.OnSeconds = on
		schedule.OffSeconds = off
		schedule.NextActionAt = nextActionAt
	}
}

func (wsm *WaterSectionManager) UpdateOnSeconds(idx waterSectionID, n timeUnit) {
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		schedule.OnSeconds = n
	}
}

func (wsm *WaterSectionManager) UpdateOffSeconds(idx waterSectionID, n timeUnit) {
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		schedule.OffSeconds = n
	}
}

func (wsm *WaterSectionManager) UpdateNextActionAt(idx waterSectionID, n timeUnit) {
	if int(idx) < len(wsm.Schedules) {
		schedule := wsm.Schedules[int(idx)]
		schedule.NextActionAt = n
	}
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
