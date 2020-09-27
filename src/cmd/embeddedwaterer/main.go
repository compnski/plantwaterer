package main

import (
	"device/avr"
	"machine"
)

func sendStats(mm *MoistureMonitor, sm *WaterSectionManager) {
	// Send all stats
	for idx, data := range mm.Data {
		println(
			idx, ": ",
			data.Samples[0].At, "=", data.Samples[0].Value, ",",
			data.Samples[1].At, "=", data.Samples[1].Value, ",",
			data.Samples[2].At, "=", data.Samples[2].Value, ",",
			data.Samples[3].At, "=", data.Samples[3].Value, ",",
			data.Samples[4].At, "=", data.Samples[4].Value, ",",
			data.Samples[5].At, "=", data.Samples[5].Value, ",",
			data.Samples[6].At, "=", data.Samples[6].Value, ",",
			data.Samples[7].At, "=", data.Samples[7].Value, ",",
			data.Samples[8].At, "=", data.Samples[8].Value,
		)
	}
	for idx := range sm.Sections {
		schedule := sm.Schedules[idx]
		println(idx, schedule.IsOn(), schedule.NextActionAt)
	}
}

func main() {
	initMachine()
	initTiming()

	cps := machine.CPUFrequency() / 256

	machine.UART0.Configure(machine.UARTConfig{BaudRate: 115200})
	println("Hello! CPS=", cps)
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	adcMulti := NewADCMultiplexer(MoistureIn, MoistureOut0, MoistureOut1, MoistureOut2)
	moistureMonitor := NewMoistureMonitor(8, adcMulti)
	_ = moistureMonitor

	sections := []*WaterSection{NewWaterSection(PumpOut0, ValveOut0),
		NewWaterSection(PumpOut0, ValveOut1),
		NewWaterSection(PumpOut0, ValveOut2)}
	//NewWaterSection(PumpOut0, ValveOut3)}
	sectionManager := NewWaterSectionManager(sections...)

	sectionManager.Update(0, 6, 60, seconds+1)
	sectionManager.Update(1, 6, 60, seconds+2)
	sectionManager.Update(2, 6, 60, seconds+3)
	//sectionManager.Update(3, 6, 60, seconds+4)

	for {
		led.High()

		led.Low()
		for {
			if t >= cps {
				t = 0
				seconds++
				sectionManager.Process(seconds)
				moistureMonitor.CheckAll(seconds)
				if seconds%5 == 0 {
					sendStats(moistureMonitor, sectionManager)
				}

				break
			}
			avr.Asm("sleep")
			//avr.Asm("nop")
			//runtime.Sleep()
		}
	}
}
