package cpu

import "fmt"

// program counter
var pc uint16

// stack pointer
var sp uint16

// accumulator
var a uint8

// flags register, lower 4 bits are unused and always 0
var f uint8
var b uint8
var c uint8
var d uint8
var e uint8
var h uint8
var l uint8

// interrupt master enable flag
var ime bool

func Reset() {
	fmt.Println("Go: cpu.Reset()")

	a = 0x01
	f = 0x80
	b = 0x00
	c = 0x13
	d = 0x00
	e = 0xD8
	h = 0x01
	l = 0x4D
	pc = 0x0100
	sp = 0xFFFE
}
