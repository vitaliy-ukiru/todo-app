package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/vitaliy-ukiru/todo-app/internal/delivery/restapi/helpers"
	"github.com/vitaliy-ukiru/todo-app/internal/delivery/restapi/helpers/response"
	domain "github.com/vitaliy-ukiru/todo-app/internal/user"
	"github.com/vitaliy-ukiru/todo-app/pkg/log"
)

type UserController struct {
	useCase domain.Usecase
	logger  log.Logger
}

func NewUserController(useCase domain.Usecase, logger log.Logger) *UserController {
	return &UserController{useCase: useCase, logger: logger}
}

func (u UserController) Create(c *fiber.Ctx) error {
	var body domain.CreateUserDTO
	if err := c.BodyParser(&body); err != nil {
		return response.Err(c, http.StatusUnprocessableEntity, "invalid credentials")
	}

	user, err := u.useCase.Create(c.Context(), body)
	if err != nil {
		return response.Err(c, 400, err.Error())
	}
	u.logger.Info("new user", log.UUID("user_id", user.ID))
	return response.Ok(c, http.StatusCreated, user)

}
func (u UserController) GetSelf(c *fiber.Ctx) error {
	claims, err := helpers.DecodeClaims(c) // token checks in middleware
	if err != nil {
		return response.Err(c, 403, "invalid access_token")
	}

	user, err := u.useCase.ByID(c.Context(), claims.ID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return response.WithError(c, 404, err)
		}
		u.logger.With(
			log.UUID("user_id", user.ID),
			log.Error(err),
		).Error("cannot get user")
		return response.Err(c, http.StatusInternalServerError, "cannot get user")
	}

	return response.Ok(c, 200, user)
}

func (u UserController) UpdatePassword(c *fiber.Ctx) error {
	claims, err := helpers.DecodeClaims(c)
	if err != nil {
		return response.Err(c, 403, "invalid access_token")
	}

	var body domain.UpdatePasswordUserDTO
	if err = c.BodyParser(&body); err != nil {
		return response.Err(c, http.StatusUnprocessableEntity, "invalid credentials")
	}
	body.UserID = claims.ID
	if err := u.useCase.UpdatePassword(c.Context(), body); err != nil {
		//TODO handling error
		return response.WithError(c, 400, err)
	}
	return response.Write(c, 200, response.JustOK)
}

func (u UserController) DeleteSelf(c *fiber.Ctx) error {
	claims, err := helpers.DecodeClaims(c)
	if err != nil {
		return response.Err(c, 403, "invalid access_token")
	}
	if err := u.useCase.Delete(c.Context(), claims.ID); err != nil {
		return response.WithError(c, 500, err)
	}
	u.logger.Info("delete user", log.UUID("user_id", claims.ID))
	return response.Write(c, 200, response.JustOK)
}
