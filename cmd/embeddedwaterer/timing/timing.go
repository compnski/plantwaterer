package timing

import (
	"device/avr"
	"machine"
	"runtime/interrupt"
)

type TimeUnit uint32

var seconds = TimeUnit(0)
var t = TimeUnit(0)
var CPS = TimeUnit(machine.CPUFrequency() / 256)

func trackTime(interrupt.Interrupt) {
	t++
}

// Initialize sets up the AVR timer interrupt to track time.
// It must be called before using Wait() or WaitForSecond()
func Initialize() {
	avr.TIMSK0.Set(avr.TIMSK0_TOIE0)
	avr.TCCR0B.Set(avr.TCCR0B_CS00)
	avr.TIFR0.Set(avr.TIFR0_TOV0)
	interrupt.New(avr.IRQ_TIMER0_OVF, trackTime)
}

// Wait waits for up to 256 ticks
func Wait(wait uint8) {
	var w = TimeUnit(wait)
	start := t
	for t-start < w {
		avr.Asm("nop")
	}
}

// WaitForSecond waits until the next second then returns the current number of seconds
// It returns a monotonically increasing count of seconds. If execution is somehow paused
// seconds will jump ahead when it returns, assuming the interrupt is still triggering.
func WaitForSecond() TimeUnit {
	for {
		if t >= CPS {
			t -= CPS
			seconds += 1
			// If execution somehow pauses, this will skip time ahead
			// If we don't wait, this function will return immeadiately until
			// seconds catches up AND will return be using incorrect values of seconds.
			if t < CPS {
				break
			}
		}
		avr.Asm("sleep")
	}
	return seconds
}
