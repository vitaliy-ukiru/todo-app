package zap

import (
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ConfigFunc func(config zapcore.EncoderConfig) zapcore.EncoderConfig

type Encoders struct {
	File    zapcore.Encoder
	Console zapcore.Encoder
}

func NewLoggerConsole(console io.Writer) Logger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC1123)
	return NewLoggerEncoders(console,
		Encoders{
			Console: zapcore.NewJSONEncoder(cfg),
		},
	)
}

func NewLogger(console io.Writer, files ...io.Writer) Logger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC1123)
	return NewLoggerEncoders(
		console,
		Encoders{
			Console: zapcore.NewConsoleEncoder(ConsoleEncoder(cfg)),
			File:    zapcore.NewJSONEncoder(cfg),
		},
		files...,
	)
}

func NewLoggerEncoders(console io.Writer, encoders Encoders, files ...io.Writer) Logger {
	//pe := zap.NewProductionEncoderConfig()
	//
	//// file
	//
	//fileEncoder := zapcore.NewJSONEncoder(pe)
	//
	//consoleCfg := consoleEncFn(pe)
	//consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)
	// //
	cores := make([]zapcore.Core, len(files)+1)

	// console
	cores[0] = zapcore.NewCore(encoders.Console,
		zapcore.AddSync(console),
		zap.DebugLevel,
	)

	// // add syncers
	for i := range files {
		cores[i+1] = zapcore.NewCore(encoders.File,
			zapcore.AddSync(files[i]),
			zap.DebugLevel,
		)
	}

	//
	return Logger{zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
	)}
}
