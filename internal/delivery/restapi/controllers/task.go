package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/vitaliy-ukiru/todo-app/internal/delivery/restapi/helpers"
	"github.com/vitaliy-ukiru/todo-app/internal/delivery/restapi/helpers/response"
	domain "github.com/vitaliy-ukiru/todo-app/internal/task"
	"github.com/vitaliy-ukiru/todo-app/pkg/log"
)

type TaskController struct {
	uc     domain.Usecase
	logger log.Logger
}

func NewTaskController(uc domain.Usecase, logger log.Logger) TaskController {
	return TaskController{uc: uc, logger: logger}
}

func (t TaskController) CreateTask(c *fiber.Ctx) error {
	var body domain.CreateTaskDTO
	if err := c.BodyParser(&body); err != nil {
		return response.Wrap(c, http.StatusUnprocessableEntity, "malformed body", err)
	}
	task, err := t.uc.Create(c.Context(), body)
	if err != nil {
		return response.Wrap(c, 500, "cannot create task", err)
	}
	return response.Ok(c, 201, task)
}

func (t TaskController) GetTask(c *fiber.Ctx) error {
	taskId, err := helpers.ParseParamUUID(c)
	if err != nil {
		return response.Wrap(c, 400, "invalid task id", err)
	}
	task, err := t.uc.FindByID(c.Context(), taskId)
	if err != nil {
		return response.Wrap(c, 400, "cannot get task", err)
	}
	return response.Ok(c, 200, task)
}

func (t TaskController) UpdateTask(c *fiber.Ctx) error {
	var body domain.UpdateFieldsDTO
	taskId, err := helpers.ParseParamUUID(c)
	if err != nil {
		return response.Wrap(c, 400, "invalid task id", err)
	}

	if err := c.BodyParser(&body); err != nil {
		return response.Wrap(c, 422, "malformed dto", err)
	}

	task, err := t.uc.UpdateTask(c.Context(), domain.UpdateTaskDTO{
		TaskID: taskId,
		Title:  body.NewTitle,
		Body:   body.NewBody,
		Status: body.NewStatus,
	})
	if err != nil {
		return response.Wrap(c, 400, "cannot update task", err)
	}
	return response.Ok(c, http.StatusAccepted, task)
}

func (t TaskController) ChangeStatus(c *fiber.Ctx) error {
	taskId, err := helpers.ParseParamUUID(c)
	if err != nil {
		return response.Wrap(c, 400, "invalid task id", err)
	}

	status, err := t.uc.ChangeStatus(c.Context(), taskId)
	if err != nil {
		return response.WithError(c, 400, err)
	}
	return response.Ok(c, http.StatusAccepted, fiber.Map{
		"status_now": status,
	})
}

func (t TaskController) DeleteTask(c *fiber.Ctx) error {
	taskId, err := helpers.ParseParamUUID(c)
	if err != nil {
		return response.Wrap(c, 400, "invalid task id", err)
	}

	if err := t.uc.Delete(c.Context(), taskId); err != nil {
		return response.WithError(c, 500, err)
	}
	return response.Write(c, 200, response.JustOK)
}

// DefaultListTasks responds array of task from default user list
func (t TaskController) DefaultListTasks(c *fiber.Ctx) error {
	claims, err := helpers.DecodeClaims(c)
	if err != nil {
		return response.Wrap(c, 403, "invalid access_token", err)
	}
	tasks, err := t.uc.MainList(c.Context(), claims.ID)
	if err != nil {
		return response.WithError(c, 500, err)
	}
	return response.Ok(c, 200, tasks)
}
