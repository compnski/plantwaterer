package main

import (
	"machine"
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

type MoistureMonitor struct {
	*ADCMultiplexer
}


func main() {
	machine.InitADC()

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	adcMulti := NewADCMultiplexer(machine.ADC0, machine.D2, machine.D3, machine.D4)
	_ = adcMulti
	

	
}
