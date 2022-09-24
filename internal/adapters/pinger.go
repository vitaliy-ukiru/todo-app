package adapters

import (
	"context"
	"time"

	"github.com/vitaliy-ukiru/todo-app/internal/delivery/restapi/controllers"
)

type JustPinger interface {
	Ping(ctx context.Context) (time.Duration, error)
}

type pingerWithName struct {
	name string
	JustPinger
}

func (p pingerWithName) Name() string {
	return p.name
}

func PingerAdapter(name string, p JustPinger) controllers.Pinger {
	return &pingerWithName{name: name, JustPinger: p}
}
