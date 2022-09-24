package controllers

import (
	"github.com/gofiber/fiber/v2"
	domain "github.com/vitaliy-ukiru/todo-app/internal/list"
	"github.com/vitaliy-ukiru/todo-app/pkg/log"
)

type ListController struct {
	uc     domain.Usecase
	logger log.Logger
}

func NewListController(uc domain.Usecase, logger log.Logger) ListController {
	return ListController{uc: uc, logger: logger}
}

func (l ListController) Create(c *fiber.Ctx) error       { return nil }
func (l ListController) AllLists(c *fiber.Ctx) error     { return nil }
func (l ListController) GetOneList(c *fiber.Ctx) error   { return nil }
func (l ListController) Update(c *fiber.Ctx) error       { return nil }
func (l ListController) Delete(c *fiber.Ctx) error       { return nil }
func (l ListController) GetInfo(c *fiber.Ctx) error      { return nil }
func (l ListController) GetListTasks(c *fiber.Ctx) error { return nil }
