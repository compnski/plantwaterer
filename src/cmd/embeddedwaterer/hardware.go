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
	MoistureOut0 Pin = machine.D6
	MoistureOut1 Pin = machine.D7
	MoistureOut2 Pin = machine.D8

	PumpOut0  Pin = machine.D2
	ValveOut0 Pin = machine.D3
	ValveOut1 Pin = machine.D13
	ValveOut2 Pin = machine.D13
	//ValveOut3 Pin = machine.D9

	//LCDSerialTx Pin = machine.D10
	DisplaySerialTX  Pin = machine.D9
	DispalySerialCLK Pin = machine.D10

	RasbPiSerialTx Pin = machine.D11
	RasbPiSerialRx Pin = machine.D12
)

func initMachine() {
	machine.InitADC()
}
