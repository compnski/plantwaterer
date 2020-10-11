package main

import (
	"github.com/compnski/plantwaterer/cmd/embeddedwaterer/timing"
	"github.com/compnski/plantwaterer/cmd/embeddedwaterer/util"
)

type MultiADC interface {
	SelectInput(n int8)
	Read() uint16
}

type MoistureSample struct {
	At    timing.TimeUnit
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

func (b *SampleBuffer) Add(now timing.TimeUnit, value uint16) {
	b.Samples[b.WriteIndex].Value = value
	b.Samples[b.WriteIndex].At = now
	b.WriteIndex = (b.WriteIndex + 1) % b.MaxSamples
}

type MoistureMonitor struct {
	ADCMultiplexer MultiADC
	NumSensors     int8
	Data           []SampleBuffer
}

const SampleBufferSize = 10
const SampleStatsToReturn = 9

func init() {
	if SampleStatsToReturn > SampleBufferSize {
		panic("SampleStatsToReturn > SampleBufferSize. Will result in error.")
	}
}

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

// SendStats prints a JSON object with historical samples for each sensor.
// Returns SampleStatsToReturn items, which MUST BE less than SampleBufferSize

func (mm *MoistureMonitor) SendStats(now timing.TimeUnit) {
	for idx, data := range mm.Data {
		print(util.MaybeLeadingComma(idx), `{"id":`, idx, `,"d":[`)
		for sampleIdx := 0; sampleIdx < SampleStatsToReturn; sampleIdx++ {
			print(util.MaybeLeadingComma(sampleIdx),
				data.Samples[sampleIdx].At, ",", data.Samples[sampleIdx].Value)
		}
		println("]}")
	}
}

// MoistureSensorCheckDelay is a number of ticks to wait between selecting a moisture sensor
// and reading its value. This can be 0 for the DirectADC, but is likely non-zero for the
// multiplexer ADC, to allow IC switching time.
const MoistureSensorCheckDelay uint8 = 50

// Check selects an input, waits some amount of ticks then reads the sensor.
func (mm *MoistureMonitor) Check(now timing.TimeUnit, n int8) {
	mm.ADCMultiplexer.SelectInput(n)
	timing.Wait(MoistureSensorCheckDelay)

	d := mm.ADCMultiplexer.Read()
	mm.Data[n].Add(now, d)
}

// CheckAll checks all moisture sensors, recording their current values
// into the SampleBuffers.
func (mm *MoistureMonitor) CheckAll(now timing.TimeUnit) {
	for idx := int8(0); idx < mm.NumSensors; idx++ {
		mm.Check(now, idx)
	}
}
