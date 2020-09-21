package main

import (	
	"machine"
)


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

