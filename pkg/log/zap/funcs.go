package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	String = zap.String
	Int64  = zap.Int64
	Error  = zap.Error
	Bool   = zap.Bool
	Any    = zap.Any
)

func ConsoleEncoder(p zapcore.EncoderConfig) zapcore.EncoderConfig {
	p.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(caller.TrimmedPath())
		encoder.AppendString("|")
	}

	p.EncodeTime = zapcore.TimeEncoderOfLayout("02/01 15:04:05") // "02/01/2006 15:04:05 |"
	p.ConsoleSeparator = " "
	p.EncodeName = func(n string, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(n)
		enc.AppendString("|")
	}

	p.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("|")
		enc.AppendString(l.CapitalString())
		enc.AppendString("|")
	}

	return p

}
