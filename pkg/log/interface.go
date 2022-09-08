package log

import "github.com/vitaliy-ukiru/todo-app/pkg/log/zap"

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Named(s string) zap.Logger
	With(fields ...Field) zap.Logger

	Sync() error
}
