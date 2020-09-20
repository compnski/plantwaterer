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

type ADCMultiplexer struct {
	// Apparently there is no garbage collection? TBD how much this static allocation is needed
	OutputPins []Pin
	InputADC ADC
}

func NewADCMultiplexer(inputPin machine.Pin, output... Pin) *ADCMultiplexer {
	adc := &ADCMultiplexer{
		InputADC: machine.ADC{inputPin},
	}
	adc.OutputPins = output
	adc.Initialize()
	return adc
}

func (adc *ADCMultiplexer) Initialize() {
	for _, outputPin := range adc.OutputPins {
		outputPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	adc.InputADC.Configure()
}

var BitMasks = [8]int8{1,2,4,8,16,32,64}

func (adc *ADCMultiplexer) SelectInput(n int8) {
	for idx, outputPin := range adc.OutputPins {
		mask := BitMasks[idx]
		if n & mask == mask {
			outputPin.High()
		} else {
			outputPin.Low()
		}
	}
}

func (adc *ADCMultiplexer) Read() uint16 {
	return adc.InputADC.Get()
}

type MultiADC interface {
	SelectInput(n int8)
	Read() uint16
}


type MoistureSample struct {
	At timeUnit
	Value uint16
}

type SampleBuffer struct {
	Samples []MoistureSample
	MaxSamples int8
	WriteIndex int8
}

func  (b *SampleBuffer) Initialize(n int8) {	
	b.Samples= make([]MoistureSample, n)
	b.MaxSamples = n
}


func (b *SampleBuffer) Add(value uint16) {
	b.Samples[b.WriteIndex].Value = value
	b.Samples[b.WriteIndex].At = timeUnit(seconds)
	println(b.WriteIndex,
		b.Samples[b.WriteIndex].Value,
		b.Samples[b.WriteIndex].At)
			
	b.WriteIndex = (b.WriteIndex + 1) % b.MaxSamples
}

type MoistureMonitor struct {
	ADCMultiplexer *ADCMultiplexer
	NumSensors int8
	Data []SampleBuffer
}
type timeUnit uint16

const SampleBufferSize = 10

func NewMoistureMonitor(numMonitors int8, m *ADCMultiplexer) *MoistureMonitor {
	mm := &MoistureMonitor{
		ADCMultiplexer: m,
		NumSensors: numMonitors,
		Data: make([]SampleBuffer, numMonitors),
	}
	for i := range mm.Data {
		mm.Data[i].Initialize(SampleBufferSize)
	}
	return mm
}

func (mm  *MoistureMonitor) Check(n int8) {
	mm.ADCMultiplexer.SelectInput(n)
	mm.Data[n].Add(mm.ADCMultiplexer.Read())
	//time.Sleep(time.Millisecond)
	//mm.Data[n].Add(mm.ADCMultiplexer.Read())
}

func (mm  *MoistureMonitor) CheckAll() {
	for idx := int8(0); idx < mm.NumSensors; idx++ {
		mm.Check(idx)
	}
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


