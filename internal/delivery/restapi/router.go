package restapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vitaliy-ukiru/todo-app/internal/delivery/restapi/controllers"
	"github.com/vitaliy-ukiru/todo-app/internal/delivery/restapi/middleware"
	"github.com/vitaliy-ukiru/todo-app/pkg/log"
)

type Router struct {
	logger log.Logger

	authMiddleware middleware.AuthMiddlewareFactory

	authController controllers.AuthController
	userController controllers.UserController
	taskController controllers.TaskController
	listController controllers.ListController

	pingController controllers.PingController
}

func NewRouter(
	logger log.Logger,
	authMiddleware middleware.AuthMiddlewareFactory,
	authController controllers.AuthController,
	userController controllers.UserController,
	taskController controllers.TaskController,
	listController controllers.ListController,
	pingController controllers.PingController,
) *Router {
	return &Router{
		logger:         logger,
		authMiddleware: authMiddleware,
		authController: authController,
		userController: userController,
		taskController: taskController,
		listController: listController,
		pingController: pingController,
	}
}

func (r Router) Bind(api fiber.Router) {
	api.Route("/auth", func(auth fiber.Router) {
		auth.Post("/login", r.authController.Login)
		auth.Get(
			"/logout",
			r.authMiddleware.Auth(),
			r.authController.Logout,
		)

	}, "auth")

	api.Route("/user", func(users fiber.Router) {
		users.Post("/", r.userController.Create)

		// with auth
		users.Route("/", func(auth fiber.Router) {
			auth.Use(r.authMiddleware.Auth())

			auth.Get("/", r.userController.GetSelf)
			auth.Delete("/", r.userController.DeleteSelf)

			auth.Put("/password", r.userController.UpdatePassword)
		})
	}, "users")

	api.Route("/task", func(tasks fiber.Router) {
		tasks.Use(r.authMiddleware.Auth())

		tasks.Post("/", r.taskController.CreateTask)
		tasks.Get("/:id", r.taskController.GetTask)
		tasks.Put("/:id", r.taskController.UpdateTask)
		tasks.Put("/:id/change-status", r.taskController.ChangeStatus)
		tasks.Delete("/:id", r.taskController.DeleteTask)

		tasks.Get("/default-list", r.taskController.DefaultListTasks)

	}, "tasks")

	api.Route("/task-list", func(lists fiber.Router) {
		lists.Use(r.authMiddleware.Auth())

		api.Post("/", r.listController.Create)
		api.Get("/", r.listController.AllLists)

		api.Get("/:id", r.listController.GetOneList)
		api.Put("/:id", r.listController.Update)
		api.Delete("/:id", r.listController.Delete)

		api.Get("/:id/info", r.listController.GetInfo)
		api.Get("/:id/tasks", r.listController.GetListTasks)

	}, "task-lists")

	api.Get("/ping", r.pingController.Ping)
}
