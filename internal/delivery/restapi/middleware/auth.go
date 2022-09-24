package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vitaliy-ukiru/todo-app/internal/auth"
	"github.com/vitaliy-ukiru/todo-app/internal/delivery/restapi/helpers"
	"github.com/vitaliy-ukiru/todo-app/internal/delivery/restapi/helpers/response"
)

type AuthMiddlewareFactory struct {
	uc *auth.JwtService
}

func NewAuthFactory(uc *auth.JwtService) AuthMiddlewareFactory {
	return AuthMiddlewareFactory{uc: uc}
}

func (a AuthMiddlewareFactory) Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := helpers.GetAccessToken(c)
		if err != nil {
			return response.Err(c, 401, "cannot get access_token")
		}

		if err := a.uc.Verify(token); err != nil {
			return response.Err(c, 401, "incorrect access_token")
		}

		return c.Next()
	}
}
