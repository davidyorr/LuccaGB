//go:build trace
// +build trace

package logger

type TraceLogger struct {
	buffer  []byte
	offset  int
	maxSize int
}

var (
	GlobalTraceLogger *TraceLogger
)

const (
	LogTypeInstruction = 0
	LogTypeMemRead     = 1
	LogTypeMemWrite    = 2
)

// Initialize on package import (or call explicitly)
func init() {
	GlobalTraceLogger = &TraceLogger{
		buffer:  make([]byte, 32*1024*1024),
		offset:  0,
		maxSize: 32 * 1024 * 1024,
	}
}

// Log instruction with optional data: [type:1][pc:2][opcode:1][data:4]
func (t *TraceLogger) LogInstruction(pc uint16, opcode byte, optionalData ...int) {
	if t.offset+8 > t.maxSize {
		return
	}

	t.buffer[t.offset] = LogTypeInstruction
	t.buffer[t.offset+1] = byte(pc >> 8)
	t.buffer[t.offset+2] = byte(pc)
	t.buffer[t.offset+3] = opcode

	// Write optional data (default to 0 if not provided)
	data := 0
	if len(optionalData) > 0 {
		data = optionalData[0]
	}

	t.buffer[t.offset+4] = byte(data >> 24)
	t.buffer[t.offset+5] = byte(data >> 16)
	t.buffer[t.offset+6] = byte(data >> 8)
	t.buffer[t.offset+7] = byte(data)
	t.offset += 8
}

// Log memory read: [type:1][addr:2][value:1]
func (t *TraceLogger) LogMemRead(addr uint16, value byte) {
	if t.offset+4 > t.maxSize {
		return
	}
	t.buffer[t.offset] = LogTypeMemRead
	t.buffer[t.offset+1] = byte(addr >> 8)
	t.buffer[t.offset+2] = byte(addr)
	t.buffer[t.offset+3] = value
	t.offset += 4
}

// Log memory write: [type:1][addr:2][value:1]
func (t *TraceLogger) LogMemWrite(addr uint16, value byte) {
	if t.offset+4 > t.maxSize {
		return
	}
	t.buffer[t.offset] = LogTypeMemWrite
	t.buffer[t.offset+1] = byte(addr >> 8)
	t.buffer[t.offset+2] = byte(addr)
	t.buffer[t.offset+3] = value
	t.offset += 4
}

func (t *TraceLogger) GetBuffer() []byte {
	return t.buffer[:t.offset]
}

func (t *TraceLogger) Reset() {
	t.offset = 0
}
