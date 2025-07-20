package timer

type Timer struct {
	// 0xFF05 timer counter
	tima uint8
	// 0xFF06 timer modulo
	tma uint8
	// 0xFF07 timer control
	tac uint8

	counter               uint16
	previousTimerBitState bool
	timaReloading         bool
	timaReloadDelay       uint8
}

func New() *Timer {
	timer := &Timer{}

	return timer
}

func (timer *Timer) Reset() {
	timer.tima = 0x00
	timer.tma = 0x00
	timer.tac = 0xF8
	timer.counter = 0xABCC
	timer.previousTimerBitState = false
}

// Step performs 1 T-cycle of work
func (timer *Timer) Step() (requestInterrupt bool) {
	// there is a 4 cycle delay before the TIMA register is reloaded with the
	// TMA register after an overflow
	timer.timaReloading = false
	if timer.timaReloadDelay > 0 {
		timer.timaReloadDelay--

		if timer.timaReloadDelay == 0 {
			timer.tima = timer.tma
			timer.timaReloading = true
			requestInterrupt = true
		}
	}

	timer.counter++
	currentTimerBitState := timer.getTimerBitState()

	if timer.isTimerEnabled() {
		if timer.timaReloadDelay == 0 {
			if timer.previousTimerBitState && !currentTimerBitState {
				timer.tima++
				// check for overflow
				if timer.tima == 0 {
					timer.tima = 0x00
					timer.timaReloadDelay = 4
				}
			}
		}
	}

	timer.previousTimerBitState = currentTimerBitState

	return requestInterrupt
}

func (timer *Timer) Read(address uint16) uint8 {
	switch address & 0x000F {
	case 0x04:
		return uint8(timer.counter >> 8)
	case 0x05:
		return timer.tima
	case 0x06:
		return timer.tma
	case 0x07:
		// upper 5 bits are unused and should always be set
		return timer.tac | 0b1111_1000
	default:
		return 0
	}
}

func (timer *Timer) Write(address uint16, value uint8) {
	switch address & 0x000F {
	case 0x04:
		timer.counter = 0
	case 0x05:
		// If you write to TIMA during the cycle that TMA is being loaded to it
		// [B], the write will be ignored and TMA value will be written to TIMA
		// instead.
		// See: https://gbdev.gg8.se/wiki/articles/Timer_Obscure_Behaviour
		if timer.timaReloading {
			timer.tima = timer.tma
			return
		}
		// During the strange cycle [A] you can prevent the IF flag from being
		// set and prevent the TIMA from being reloaded from TMA by writing a
		// value to TIMA. That new value will be the one that stays in the TIMA
		// register after the instruction. Writing to DIV, TAC or other
		// registers won't prevent the IF flag from being set or TIMA from being
		// reloaded.
		// See: https://gbdev.gg8.se/wiki/articles/Timer_Obscure_Behaviour
		if timer.timaReloadDelay != 0 {
			timer.timaReloadDelay = 0
		}
		timer.tima = value
	case 0x06:
		// If TMA is written the same cycle it is loaded to TIMA [B], TIMA is
		// also loaded with that value.
		// See: https://gbdev.gg8.se/wiki/articles/Timer_Obscure_Behaviour
		if timer.timaReloading {
			timer.tima = value
		}
		timer.tma = value
	case 0x07:
		// upper 5 bits are unused
		timer.tac = (value & 0b0000_0111)
	}
}

func (timer *Timer) getTimerBitState() bool {
	// the bit index to check on the counter
	var bitIndex uint8 = 9

	switch timer.tac & 0b11 {
	case 0b00: // 4096 Hz
		bitIndex = 9
	case 0b01: // 262144 Hz
		bitIndex = 3
	case 0b10: // 65536 Hz
		bitIndex = 5
	case 0b11: // 16384 Hz
		bitIndex = 7
	}

	mask := uint16(1 << bitIndex)

	return (timer.counter & mask) != 0
}

func (timer *Timer) isTimerEnabled() bool {
	// bit 2 -> timer enabled
	return (timer.tac & 0b0000_0100) != 0
}
