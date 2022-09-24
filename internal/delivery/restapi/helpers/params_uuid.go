package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type params struct {
	ID uuid.UUID `params:"id"`
}

func ParseParamUUID(c *fiber.Ctx) (uuid.UUID, error) {
	var p params
	if err := c.ParamsParser(&p); err != nil {
		return uuid.Nil, err
	}
	return p.ID, nil
}
