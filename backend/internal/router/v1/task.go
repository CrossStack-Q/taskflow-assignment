package v1

import (
	"taskflow/internal/handler"
	"taskflow/internal/middleware"

	"github.com/labstack/echo/v4"
)

func registerTaskRoutes(r *echo.Group, h *handler.TaskHandler, ch *handler.CommentHandler, auth *middleware.AuthMiddleware) {

	tasks := r.Group("/tasks")
	tasks.Use(auth.RequireAuth)

	tasks.POST("", h.CreateTask)
	tasks.GET("", h.GetTasks)
	tasks.GET("/stats", h.GetTaskStats)

	dynamicTask := tasks.Group("/:id")
	dynamicTask.GET("", h.GetTaskByID)
	dynamicTask.PATCH("", h.UpdateTask)
	dynamicTask.DELETE("", h.DeleteTask)

	taskComments := dynamicTask.Group("/comments")
	taskComments.POST("", ch.AddComment)
	taskComments.GET("", ch.GetCommentsByTaskID)

}
