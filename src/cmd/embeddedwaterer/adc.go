package main

import (
	"machine"
)

type ADCDirect struct {
	InputADCs   []ADC
	selectedADC int8
}

func NewADCDirect(inputADCs ...machine.Pin) *ADCDirect {
	adc := &ADCDirect{
		InputADCs: make([]ADC, len(inputADCs)),
	}
	for idx, adcPin := range inputADCs {
		adc.InputADCs[idx] = &machine.ADC{adcPin}
		adc.InputADCs[idx].Configure()
	}
	return adc
}

func (adc *ADCDirect) Inputs() int8 {
	return int8(len(adc.InputADCs))
}

func (adc *ADCDirect) SelectInput(n int8) {
	adc.selectedADC = n

}

func (adc *ADCDirect) Read() uint16 {
	if int(adc.selectedADC) < len(adc.InputADCs) {
		return adc.InputADCs[adc.selectedADC].Get()
	}
	return 0xFFFF
}

type ADCMultiplexer struct {
	// Apparently there is no garbage collection? TBD how much this static allocation is needed
	OutputPins []Pin
	InputADC   ADC
}

func NewADCMultiplexer(inputPin machine.Pin, output ...Pin) *ADCMultiplexer {
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

var BitMasks = [8]int8{1, 2, 4, 8, 16, 32, 64}

func (adc *ADCMultiplexer) SelectInput(n int8) {
	for idx, outputPin := range adc.OutputPins {
		mask := BitMasks[idx]
		if n&mask == mask {
			outputPin.High()
		} else {
			outputPin.Low()
		}
	}
}

func (adc *ADCMultiplexer) Read() uint16 {
	return adc.InputADC.Get()
}
