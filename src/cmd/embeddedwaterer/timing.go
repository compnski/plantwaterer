package main

import (
	"device/avr"
	"runtime/interrupt"
)

func now() timeUnit {
	return seconds
}

var seconds = timeUnit(0)
var t = uint32(0)

type timeUnit uint32

func trackTime(interrupt.Interrupt) {
	t++
}

func initTiming() {
	avr.TIMSK0.Set(avr.TIMSK0_TOIE0)
	avr.TCCR0B.Set(avr.TCCR0B_CS00)
	avr.TIFR0.Set(avr.TIFR0_TOV0)
	interrupt.New(avr.IRQ_TIMER0_OVF, trackTime)
}
