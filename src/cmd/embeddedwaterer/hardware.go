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

var (
	MoistureIn       = machine.ADC0
	MoistureOut0 Pin = machine.D2
	MoistureOut1 Pin = machine.D3
	MoistureOut2 Pin = machine.D4

	PumpOut0  Pin = machine.D5
	ValveOut0 Pin = machine.D6
	ValveOut1 Pin = machine.D7
	ValveOut2 Pin = machine.D8
	ValveOut3 Pin = machine.D9

	LCDSerialTx Pin = machine.D10

	RasbPiSerialTx Pin = machine.D11
	RasbPiSerialRx Pin = machine.D12
)

func initMachine() {
	machine.InitADC()
}
