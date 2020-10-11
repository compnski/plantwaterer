package hardware

import "machine"

type WaterControlDevice struct {
	Pin Pin
}

func (d *WaterControlDevice) On() {
	d.Pin.High()
}

func (d *WaterControlDevice) Off() {
	d.Pin.Low()
}

func NewWaterSection(pins ...Pin) *WaterSection {
	d := &WaterSection{
		Devices: make([]*WaterControlDevice, len(pins)),
	}
	for idx, pin := range pins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		d.Devices[idx] = &WaterControlDevice{pin}
	}
	return d
}

type WaterSection struct {
	Devices []*WaterControlDevice
}

func (s *WaterSection) On() {
	for _, device := range s.Devices {
		device.On()
	}
}

func (s *WaterSection) Off() {
	for _, device := range s.Devices {
		device.Off()
	}
}
