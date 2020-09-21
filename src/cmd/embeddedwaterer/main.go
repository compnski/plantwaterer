package main

import (	
	"machine"
	"device/avr"
	"runtime/interrupt"
)

type ADC interface {
	Get() uint16
	Configure()
}

type Pin interface {
	High()
	Low()
	Configure(machine.PinConfig)
}


var t = uint32(0)

func trackTime(interrupt.Interrupt) {
	t++
	//	println(t)
}

var seconds = uint16(0)


func sendStats(mm *MoistureMonitor) {
	// Send all stats
	for _, data := range mm.Data {
		println(
			data.Samples[0].At,"=",data.Samples[0].Value,", ",
			data.Samples[1].At,"=",data.Samples[1].Value,", ",
			data.Samples[2].At,"=",data.Samples[2].Value,", ",
			data.Samples[3].At,"=",data.Samples[3].Value,", ",
			data.Samples[4].At,"=",data.Samples[4].Value,", ",
			data.Samples[5].At,"=",data.Samples[5].Value,", ",
			data.Samples[6].At,"=",data.Samples[6].Value,", ",
			data.Samples[7].At,"=",data.Samples[7].Value,", ",
			data.Samples[8].At,"=",data.Samples[8].Value,", ",
			data.Samples[9].At,"=",data.Samples[9].Value,", ",
		)
	}
	
}

func PutUint16(b []byte, v uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}


type WaterControlDevice struct {
	Pin Pin
}

func (d *WaterControlDevice) Open() {
	d.Pin.High()
}

func (d *WaterControlDevice) Close() {
	d.Pin.Low()
}

type WaterSection struct {
	Devices []*WaterControlDevice
	WaterSectionId int8
}

func (s *WaterSection) Open() {
	for _, device := range s.Devices {
		device.Open()
	}
}

func (s *WaterSection) Close() {
	for _, device := range s.Devices {
		device.Close()
	}
}

type SectionSchedule struct {
	OnSeconds, OffSeconds uint16
	NextOnSecond, NextOffSecond uint16
	WaterSectionId int8
	State bool
}

type WaterSectionManager struct {
	Sections map[int8]WaterSection
	Schedules []SectionSchedule
}





func main() {
	machine.InitADC()

	avr.TIMSK0.Set(avr.TIMSK0_TOIE0)
	avr.TCCR0B.Set(avr.TCCR0B_CS00)
	avr.TIFR0.Set(avr.TIFR0_TOV0)
	interrupt.New(avr.IRQ_TIMER0_OVF, trackTime)
	
	cps := machine.CPUFrequency() / 256

	machine.UART0.Configure(machine.UARTConfig{BaudRate:115200})
	println("Hello! CPS=",cps)
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	adcMulti := NewADCMultiplexer(machine.ADC0, machine.D2, machine.D3, machine.D4)
	moistureMonitor := NewMoistureMonitor(8, adcMulti)
	_ = moistureMonitor
	
	for ;; {
		led.High()

		led.Low()
		println("z")
		//time.Sleep(time.Second)
		for ;; {
			//println(t);
			if t >= cps {
				moistureMonitor.CheckAll()
				sendStats(moistureMonitor)
				t = 0
				seconds++
				break
			}
			avr.Asm("nop")
			//runtime.Sleep()
		}
	}
	
	
}


