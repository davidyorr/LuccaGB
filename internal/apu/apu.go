package apu

type APU struct {
	// ======================================
	// ====== Global Control Registers ======
	// ======================================

	// 0xFF26 - NR52: Audio master control
	nr52 uint8
	// 0xFF25 - NR51: Sound panning
	nr51 uint8
	// 0xFF24 - NR50: Master volume & VIN panning
	nr50 uint8

	// =======================================================
	// ====== Sound Channel 1 — Pulse with period sweep ======
	// =======================================================

	// 0xFF10 - NR10: Channel 1 sweep
	nr10 uint8
	// 0xFF11 - NR11: Channel 1 length timer & duty cycle
	nr11 uint8
	// 0xFF12 - NR12: Channel 1 volume & envelope
	nr12 uint8
	// 0xFF13 - NR13: Channel 1 period low [write-only]
	nr13 uint8
	// 0xFF14 - NR14: Channel 1 period high & control
	nr14 uint8

	// =====================================
	// ====== Sound Channel 2 — Pulse ======
	// =====================================

	// 0xFF16 - NR21: Channel 2 length timer & duty cycle
	nr21 uint8
	// 0xFF17 - NR22: Channel 2 volume & envelope
	nr22 uint8
	// 0xFF18 - NR23: Channel 2 period low [write-only]
	nr23 uint8
	// 0xFF19 - NR24: Channel 2 period high & control
	nr24 uint8

	// ===========================================
	// ====== Sound Channel 3 — Wave output ======
	// ===========================================

	// 0xFF1A - NR30: Channel 3 DAC enable
	nr30 uint8
	// 0xFF1B - NR31: Channel 3 length timer [write-only]
	nr31 uint8
	// 0xFF1C - NR32: Channel 3 output level
	nr32 uint8
	// 0xFF1D - NR33: Channel 3 period low [write-only]
	nr33 uint8
	// 0xFF1E - NR34: Channel 3 period high & control
	nr34 uint8

	// 0xFF30-0xFF3F - Wave pattern RAM
	waveRam [16]uint8

	// =====================================
	// ====== Sound Channel 4 — Noise ======
	// =====================================

	// 0xFF20 - NR41: Channel 4 length timer [write-only]
	nr41 uint8
	// 0xFF21 - NR42: Channel 4 volume & envelope
	nr42 uint8
	// 0xFF22 - NR43: Channel 4 frequency & randomness
	nr43 uint8
	// 0xFF23 - NR44: Channel 4 control
	nr44 uint8
}

func New() *APU {
	apu := &APU{}

	apu.Reset()

	return apu
}

func (apu *APU) Reset() {
	apu.nr10 = 0x80
	apu.nr11 = 0xBF
	apu.nr12 = 0xF3
	apu.nr13 = 0xFF
	apu.nr14 = 0xBF
	apu.nr21 = 0x3F
	apu.nr22 = 0x00
	apu.nr23 = 0xFF
	apu.nr24 = 0xBF
	apu.nr30 = 0x7F
	apu.nr31 = 0xFF
	apu.nr32 = 0x9F
	apu.nr33 = 0xFF
	apu.nr34 = 0xBF
	apu.nr41 = 0xFF
	apu.nr42 = 0x00
	apu.nr43 = 0x00
	apu.nr44 = 0xBF
	apu.nr50 = 0x77
	apu.nr51 = 0xF3
	apu.nr52 = 0xF1
}

// Perform 1 T-cycle of work
func (apu *APU) Step() {
}

func (apu *APU) Read(address uint16) uint8 {
	switch {
	case address == 0xFF10:
		return apu.nr10 | 0b1000_0000
	case address == 0xFF11:
		return apu.nr11
	case address == 0xFF12:
		return apu.nr12
	case address == 0xFF13:
		return apu.nr13
	case address == 0xFF14:
		return apu.nr14
	case address == 0xFF16:
		return apu.nr21
	case address == 0xFF17:
		return apu.nr22
	case address == 0xFF18:
		return apu.nr23
	case address == 0xFF19:
		return apu.nr24
	case address == 0xFF1A:
		return apu.nr30 | 0b0111_1111
	case address == 0xFF1B:
		return apu.nr31 | 0b1000_0001
	case address == 0xFF1C:
		return apu.nr32 | 0b1001_1111
	case address == 0xFF1D:
		return apu.nr33
	case address == 0xFF1E:
		return apu.nr34
	case address == 0xFF20:
		return apu.nr41 | 0b1100_0000
	case address == 0xFF21:
		return apu.nr42
	case address == 0xFF22:
		return apu.nr43
	case address == 0xFF23:
		return apu.nr44 | 0b0011_1111
	case address == 0xFF24:
		return apu.nr50
	case address == 0xFF25:
		return apu.nr51
	case address == 0xFF26:
		return apu.nr52 | 0b0111_0000
	case address >= 0xFF30 && address <= 0xFF3F:
		return apu.waveRam[address-0xFF30]
	}

	return 0xFF
}

func (apu *APU) Write(address uint16, value uint8) {

}
