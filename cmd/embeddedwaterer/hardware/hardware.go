package hardware

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
	MoistureIn0 = machine.ADC0
	MoistureIn1 = machine.ADC1
	MoistureIn2 = machine.ADC2
	MoistureIn3 = machine.ADC3
	MoistureIn4 = machine.ADC4

	// MoistureOut0 Pin = machine.D6
	// MoistureOut1 Pin = machine.D7
	// MoistureOut2 Pin = machine.D8

	PumpOut0 Pin = machine.D4
	PumpOut1 Pin = machine.D5

	SwitchIn0 Pin = machine.D8
	SwitchIn1 Pin = machine.D9
	LedOut0   Pin = machine.D10
	SerialOut Pin = machine.D11

	// ValveOut0 Pin = machine.D3
	// ValveOut1 Pin = machine.D13
	// ValveOut2 Pin = machine.D13
	// ValveOut3 Pin = machine.D9

	//DisplaySerialTX  Pin = machine.D9
	//DispalySerialCLK Pin = machine.D10

	//RasbPiSerialTx machine.Pin = machine.D11
	//RasbPiSerialRx machine.Pin = machine.D12
)

func Initialize() {
	machine.InitADC()
	machine.UART0.Configure(machine.UARTConfig{BaudRate: 115200})
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
}
