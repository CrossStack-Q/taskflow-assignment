package v1

import (
	"taskflow/internal/handler"
	"taskflow/internal/middleware"

	"github.com/labstack/echo/v4"
)

func registerProjectRoutes(r *echo.Group, h *handler.ProjectHandler, auth *middleware.AuthMiddleware) {

	projects := r.Group("/projects")
	projects.Use(auth.RequireAuth)
	projects.POST("", h.CreateProject)
	projects.GET("", h.GetProjects)

	dynamicProject := projects.Group("/:id")
	dynamicProject.PATCH("", h.UpdateProject)
	dynamicProject.DELETE("", h.DeleteProject)
}
