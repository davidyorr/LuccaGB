package interrupt

type Interrupt uint8

const (
	VBlankInterrupt Interrupt = 0b0000_0001
	LcdInterrupt    Interrupt = 0b0000_0010
	TimerInterrupt  Interrupt = 0b0000_0100
	SerialInterrupt Interrupt = 0b0000_1000
	JoypadInterrupt Interrupt = 0b0001_0000
)
