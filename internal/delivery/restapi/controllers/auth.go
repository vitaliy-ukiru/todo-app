package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/vitaliy-ukiru/todo-app/internal/auth"
	"github.com/vitaliy-ukiru/todo-app/internal/delivery/restapi/helpers/response"
)

type AuthController struct {
	authService *auth.JwtService
}

func NewAuthController(authService *auth.JwtService) *AuthController {
	return &AuthController{authService: authService}
}

func (a AuthController) Login(c *fiber.Ctx) error {
	var body auth.CredentialsDTO
	if err := c.BodyParser(&body); err != nil {
		return response.Err(c, 422, "malformed credentials")
	}

	loginResult, err := a.authService.Login(c.Context(), body)
	if err != nil {
		return response.Err(c, 401, "invalid email or password")
	}

	//allowedCookie := !c.Context().QueryArgs().Has("disable_cookie_jwt")
	//if allowedCookie {
	//	c.Cookie(&fiber.Cookie{
	//		Name:     "access_token",
	//		Value:    token,
	//		Path:     "/",
	//		Expires:  expires,
	//		Secure:   true,
	//		HTTPOnly: true,
	//	})
	//}

	return response.Ok(c, 200, loginResult)

}

func (a AuthController) Logout(c *fiber.Ctx) error {
	//todo refresh tokens...
	return response.Err(c, http.StatusNotImplemented, "delete access_token from storage :/")
}
