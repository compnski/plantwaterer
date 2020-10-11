package main

import (
	"machine"

	"github.com/compnski/plantwaterer/cmd/embeddedwaterer/hardware"
	"github.com/compnski/plantwaterer/cmd/embeddedwaterer/remote"
	"github.com/compnski/plantwaterer/cmd/embeddedwaterer/schedule"
	"github.com/compnski/plantwaterer/cmd/embeddedwaterer/timing"
)

const StatsIntervalSeconds = 5

func CreateWaterSectionManager() *schedule.WaterSectionManager {
	sections := []*hardware.WaterSection{hardware.NewWaterSection(hardware.PumpOut0, hardware.PumpOut1)}
	sectionManager := schedule.NewWaterSectionManager(sections...)

	sectionManager.Update(0, 50, 3600*72, 0)
	return sectionManager
}

func CreateMoistureMonitor() *MoistureMonitor {
	adcDirect := hardware.NewADCDirect(hardware.MoistureIn0, hardware.MoistureIn1, hardware.MoistureIn2, hardware.MoistureIn3)
	moistureMonitor := NewMoistureMonitor(adcDirect.Inputs(), adcDirect)

	return moistureMonitor
}

type Waterer struct {
	MoistureMonitor  *MoistureMonitor
	SectionManager   *schedule.WaterSectionManager
	RemoteController *remote.RemoteController
}

func CreateWaterer() *Waterer {
	sm := CreateWaterSectionManager()
	return &Waterer{
		MoistureMonitor:  CreateMoistureMonitor(),
		SectionManager:   sm,
		RemoteController: remote.NewRemoteController(sm),
	}
}

func main() {
	var now timing.TimeUnit = 0
	hardware.Initialize()
	timing.Initialize()
	waterer := CreateWaterer()
	println(`{"hi":"!", "CPS":`, timing.CPS, `}`)
	waterer.SendStats(now)
	for {
		now = timing.WaitForSecond()
		machine.LED.High()
		waterer.Tick(now)
		machine.LED.Low()
	}
}

func (w *Waterer) Tick(now timing.TimeUnit) {
	if serialData := w.checkUartRx(); serialData != nil {
		w.RemoteController.ProcessBytes(now, serialData)
	}
	w.SectionManager.Process(now)
	w.MoistureMonitor.CheckAll(now)
	if now%StatsIntervalSeconds == 0 {
		w.SendStats(now)
	}

}

// serialData is static because we have no garbage collection and
// this data likely escapes the stack.
var serialData = make([]byte, 10)

func (w *Waterer) checkUartRx() []byte {
	if machine.UART0.Buffered() > 0 {
		if _, err := machine.UART0.Read(serialData); err == nil {
			return serialData
		} else {
			println(err.Error())
		}
	}
	return nil
}

// SendStats prints all stats to stdout which is sent out over the UART
func (w *Waterer) SendStats(now timing.TimeUnit) {
	println(`{"now":`, now, `,"sensors":`)
	w.MoistureMonitor.SendStats(now)
	println(`,"sections":`)
	w.SectionManager.SendStats(now)
	println("}")
}
