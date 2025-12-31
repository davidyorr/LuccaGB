//go:build !trace
// +build !trace

package logger

type TraceLogger struct{}

var (
	GlobalTraceLogger *TraceLogger
)

func (t *TraceLogger) Enable()                                                    {}
func (t *TraceLogger) Disable()                                                   {}
func (t *TraceLogger) LogInstruction(pc uint16, opcode byte, optionalData ...int) {}
func (t *TraceLogger) LogMemRead(addr uint16, value byte)                         {}
func (t *TraceLogger) LogMemWrite(addr uint16, value byte)                        {}
func (t *TraceLogger) GetBuffer() []byte                                          { return nil }
func (t *TraceLogger) Reset()                                                     {}
