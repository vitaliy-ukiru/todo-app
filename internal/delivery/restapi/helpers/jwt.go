package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/vitaliy-ukiru/todo-app/internal/auth"
	"github.com/vitaliy-ukiru/todo-app/pkg/jwt"
)

func GetAccessToken(c *fiber.Ctx) (string, error) {
	//token := utils.CopyString(c.Cookies(jwt.CookieAccessToken))
	//if token != "" {
	//	return token, nil
	//}

	header := utils.CopyString(c.Get(jwt.HeaderAuthorization))
	return jwt.GetTokenFromHeader(header)
}

func DecodeClaims(c *fiber.Ctx) (*auth.Claims, error) {
	token, err := GetAccessToken(c)
	if err != nil {
		return nil, err
	}
	claims := new(auth.Claims)
	return claims, jwt.UnverifiedClaims(token, claims)
}
