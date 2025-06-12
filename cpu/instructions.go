package cpu

type instruction struct {
	mnemonic      string
	operandLength uint8
	execute       func()
}

var instructions = [256]instruction{
	{"NOP", 0, nop},
}
