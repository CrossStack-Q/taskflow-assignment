package handler

import (
	"net/http"

	"taskflow/internal/middleware"
	"taskflow/internal/model"
	"taskflow/internal/model/task"
	"taskflow/internal/server"
	"taskflow/internal/service"

	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	Handler
	taskService *service.TaskService
}

func NewTaskHandler(s *server.Server, taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		Handler:     NewHandler(s),
		taskService: taskService,
	}
}

func (h *TaskHandler) CreateTask(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *task.CreateTaskPayload) (*task.Task, error) {
			userID := middleware.GetUserID(c)
			return h.taskService.CreateTask(c, userID, payload)
		},
		http.StatusCreated,
		&task.CreateTaskPayload{},
	)(c)
}

func (h *TaskHandler) GetTaskByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *task.GetTaskByIDPayload) (*task.PopulatedTask, error) {
			userID := middleware.GetUserID(c)
			return h.taskService.GetTaskByID(c, userID, payload.ID)
		},
		http.StatusOK,
		&task.GetTaskByIDPayload{},
	)(c)
}

func (h *TaskHandler) GetTasks(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *task.GetTasksQuery) (*model.PaginatedResponse[task.PopulatedTask], error) {
			userID := middleware.GetUserID(c)
			// userID := "user_3CFMzCaRFQNzy2q4byc3JcUaQJL"
			return h.taskService.GetTasks(c, userID, query)
		},
		http.StatusOK,
		&task.GetTasksQuery{},
	)(c)
}

func (h *TaskHandler) UpdateTask(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *task.UpdateTaskPayload) (*task.Task, error) {
			userID := middleware.GetUserID(c)
			return h.taskService.UpdateTask(c, userID, payload)
		},
		http.StatusOK,
		&task.UpdateTaskPayload{},
	)(c)
}

func (h *TaskHandler) DeleteTask(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *task.DeleteTaskPayload) error {
			userID := middleware.GetUserID(c)
			return h.taskService.DeleteTask(c, userID, payload.ID)
		},
		http.StatusNoContent,
		&task.DeleteTaskPayload{},
	)(c)
}

func (h *TaskHandler) GetTaskStats(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *task.GetTaskStatsPayload) (*task.TaskStats, error) {
			userID := middleware.GetUserID(c)
			return h.taskService.GetTaskStats(c, userID)
		},
		http.StatusOK,
		&task.GetTaskStatsPayload{},
	)(c)
}
