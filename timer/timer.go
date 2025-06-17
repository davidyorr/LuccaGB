package timer

import "fmt"

type Timer struct {
	// 0xFF04 divider register
	div uint8
	// 0xFF05 timer counter
	tima uint8
	// 0xFF06 timer modulo
	tma uint8
	// 0xFF07 timer control
	tac uint8

	divCounter uint16
}

func New() *Timer {
	timer := &Timer{}

	return timer
}

func (timer *Timer) Reset() {
	timer.div = 0x18
	timer.divCounter = 0
	timer.tima = 0x00
	timer.tma = 0x00
	timer.tac = 0xF8
}

// there are 4 T-cycles per M-cycle, so T-cycles / M-cycles
const threshold = 256 / 4

func (timer *Timer) Step(cycles uint8) {
	timer.divCounter += uint16(cycles)
	fmt.Println("STEPPING ------------------------------------------------------------", timer.divCounter)

	if timer.divCounter >= threshold {
		timer.div++
		timer.divCounter -= threshold
	}
}

func (timer *Timer) Read(address uint16) uint8 {
	switch address & 0x000F {
	case 0x04:
		return timer.div
	case 0x05:
		return timer.tima
	case 0x06:
		return timer.tima
	case 0x07:
		return timer.tac
	default:
		return 0
	}
}

func (timer *Timer) Write(address uint16, value uint8) {
	switch address & 0x000F {
	case 0x04:
		timer.div = 0
		timer.divCounter = 0
	case 0x05:
		// timer.tima
	case 0x06:
		// timer.tima
	case 0x07:
		// timer.tac
	}
}
