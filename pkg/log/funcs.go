package log

import (
	"github.com/google/uuid"
	"github.com/vitaliy-ukiru/todo-app/pkg/log/zap"
)

func Bool(key string, value bool) Field {
	return zap.Bool(key, value)
}

func Int64(key string, value int64) Field {
	return zap.Int64(key, value)
}

func Any(key string, value any) Field {
	return zap.Any(key, value)
}

func Error(value error) Field {
	return zap.Error(value)
}

func UUID(key string, value uuid.UUID) Field {
	return zap.String(key, value.String())
}

func String(key string, value string) Field {
	return zap.String(key, value)
}
