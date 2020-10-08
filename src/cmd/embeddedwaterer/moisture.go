package main

import (
	"device/avr"
)

type MultiADC interface {
	SelectInput(n int8)
	Read() uint16
}

type MoistureSample struct {
	At    timeUnit
	Value uint16
}

type SampleBuffer struct {
	Samples    []MoistureSample
	MaxSamples int8
	WriteIndex int8
}

func (b *SampleBuffer) Initialize(n int8) {
	b.Samples = make([]MoistureSample, n)
	b.MaxSamples = n
}

func (b *SampleBuffer) Add(now timeUnit, value uint16) {
	b.Samples[b.WriteIndex].Value = value
	b.Samples[b.WriteIndex].At = now
	//println(b.Samples[b.WriteIndex].Value,
	//	b.Samples[b.WriteIndex].At)

	b.WriteIndex = (b.WriteIndex + 1) % b.MaxSamples
}

type MoistureMonitor struct {
	ADCMultiplexer MultiADC
	NumSensors     int8
	Data           []SampleBuffer
}

const SampleBufferSize = 10

func NewMoistureMonitor(numMonitors int8, m MultiADC) *MoistureMonitor {
	mm := &MoistureMonitor{
		ADCMultiplexer: m,
		NumSensors:     numMonitors,
		Data:           make([]SampleBuffer, numMonitors),
	}
	for i := range mm.Data {
		mm.Data[i].Initialize(SampleBufferSize)
	}
	return mm
}

func (mm *MoistureMonitor) Check(now timeUnit, n int8) {

	mm.ADCMultiplexer.SelectInput(n)
	start := t
	for t-start < 100 {
		avr.Asm("nop")
	}

	d := mm.ADCMultiplexer.Read()
	// avr.Asm("nop")
	// d += mm.ADCMultiplexer.Read() / 4
	// avr.Asm("nop")
	// d += mm.ADCMultiplexer.Read() / 4
	// avr.Asm("nop")
	// d += mm.ADCMultiplexer.Read() / 4
	mm.Data[n].Add(now, d)
	//time.Sleep(time.Millisecond)
	//mm.Data[n].Add(mm.ADCMultiplexer.Read())
}

func (mm *MoistureMonitor) CheckAll(now timeUnit) {
	for idx := int8(0); idx < mm.NumSensors; idx++ {
		mm.Check(now, idx)
	}
}
