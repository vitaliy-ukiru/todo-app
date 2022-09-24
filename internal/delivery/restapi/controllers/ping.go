package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vitaliy-ukiru/todo-app/internal/delivery/restapi/helpers/response"
)

type Pinger interface {
	Ping(ctx context.Context) (time.Duration, error)
	Name() string
}

type PingController struct {
	services []Pinger
}

func NewPingController(services ...Pinger) PingController {
	return PingController{services: services}
}

type pingResult struct {
	Name string `json:"name"`
	Time string `json:"time"`
	Ms   int64  `json:"ms"`
}

func (p PingController) Ping(c *fiber.Ctx) error {
	results := make([]pingResult, 0, len(p.services))
	var sum int64
	for _, service := range p.services {
		ping, err := service.Ping(c.Context())
		if err != nil {
			return response.Wrap(c, 500, service.Name()+": cannot ping", err)
		}
		results = append(results, pingResult{
			Name: service.Name(),
			Time: ping.String(),
			Ms:   ping.Milliseconds(),
		})
		sum += ping.Milliseconds()
	}
	avg := float64(sum) / float64(len(results))
	return response.Ok(c, 200, fiber.Map{
		"avg":      avg,
		"services": results,
	})
}
