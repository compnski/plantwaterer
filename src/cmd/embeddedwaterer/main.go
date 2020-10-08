package main

import (
	"device/avr"
	"machine"
)

func sendStats(mm *MoistureMonitor, sm *WaterSectionManager) {
	// Send all stats
	println(`{"now":`, now(), `,"sensors":[`)
	for idx, data := range mm.Data {
		var maybeLeadingComma string
		if idx > 0 {
			maybeLeadingComma = ","
		}
		println(maybeLeadingComma,
			`{"id":`, idx, `,"d":[`,
			data.Samples[0].At, ",", data.Samples[0].Value, ",",
			data.Samples[1].At, ",", data.Samples[1].Value, ",",
			data.Samples[2].At, ",", data.Samples[2].Value, ",",
			data.Samples[3].At, ",", data.Samples[3].Value, ",",
			data.Samples[4].At, ",", data.Samples[4].Value, ",",
			data.Samples[5].At, ",", data.Samples[5].Value, ",",
			data.Samples[6].At, ",", data.Samples[6].Value, ",",
			data.Samples[7].At, ",", data.Samples[7].Value, ",",
			data.Samples[8].At, ",", data.Samples[8].Value, "]}",
		)
	}
	println(`],"sections":[`)
	for idx := range sm.Sections {
		var maybeLeadingComma string
		if idx > 0 {
			maybeLeadingComma = ","
		}
		schedule := sm.Schedules[idx]
		println(maybeLeadingComma, `{"id":`, idx, ",",
			`"on":`, schedule.IsOn(), ",",
			`"next":`, schedule.NextActionAt, ",",
			`"last":`, schedule.LastActionAt, ",",
			`"onTime":`, schedule.OnSeconds, ",",
			`"offTime":`, schedule.OffSeconds, ",",
			`"onAcc":`, schedule.OnAccum, ",",
			`"offAcc":`, schedule.OffAccum,
			`}`,
		)
	}
	println("]}")
}

var CmdPrefix = []byte{'C', 'M', 'D'}
var serialData = make([]byte, 10)

func main() {
	initMachine()
	initTiming()

	cps := machine.CPUFrequency() / 256

	//RasbPiSerialTx.Configure(machine.PinConfig{Mode: machine.PinOutput})
	//rasbUART := machine.UART{Buffer: machine.NewRingBuffer()}
	//rasbUART.Configure(machine.UARTConfig{BaudRate: 9600, TX: RasbPiSerialTx, RX: RasbPiSerialRx})

	machine.UART0.Configure(machine.UARTConfig{BaudRate: 115200})
	//machine.UART0.Buffer = machine.NewRingBuffer()
	println(`{"hi":"!", "CPS":`, cps, `}`)
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	//adcMulti := NewADCMultiplexer(MoistureIn, MoistureOut0, MoistureOut1, MoistureOut2)
	adcDirect := NewADCDirect(MoistureIn0, MoistureIn1, MoistureIn2, MoistureIn3)
	moistureMonitor := NewMoistureMonitor(adcDirect.Inputs(), adcDirect)
	_ = moistureMonitor

	sections := []*WaterSection{NewWaterSection(PumpOut0, PumpOut1)}
	sectionManager := NewWaterSectionManager(sections...)

	sectionManager.Update(0, 50, 3600*72, now()+1)
	controller := &Controller{wsm: sectionManager}

	for {
		if machine.UART0.Buffered() > 0 {
			_, err := machine.UART0.Read(serialData)
			if err == nil {
				controller.ProcessBytes(serialData)
			} else {
				println(err.Error())
			}
		}

		led.High()
		for {
			if t >= cps {
				t = 0
				seconds += 1
				sectionManager.Process(now())
				moistureMonitor.CheckAll(now())
				if now()%5 == 0 {
					//println(machine.UART0.Buffer.Used())
					sendStats(moistureMonitor, sectionManager)
				}
				break
			}
			led.Low()
			avr.Asm("sleep")
			//avr.Asm("nop")
			//runtime.Sleep()
		}
	}
}
