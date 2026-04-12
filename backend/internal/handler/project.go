package handler

import (
	"net/http"

	"taskflow/internal/middleware"
	"taskflow/internal/model"
	"taskflow/internal/model/project"
	"taskflow/internal/server"
	"taskflow/internal/service"

	"github.com/labstack/echo/v4"
)

type ProjectHandler struct {
	Handler
	projectService *service.ProjectService
}

func NewProjectHandler(s *server.Server, projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		Handler:        NewHandler(s),
		projectService: projectService,
	}
}

func (h *ProjectHandler) CreateProject(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *project.CreateProjectPayload) (*project.Project, error) {
			userID := middleware.GetUserID(c)
			return h.projectService.CreateProject(c, userID, payload)
		},
		http.StatusCreated,
		&project.CreateProjectPayload{},
	)(c)
}

func (h *ProjectHandler) GetProjects(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *project.GetProjectsQuery) (
			*model.PaginatedResponse[project.Project], error,
		) {
			userID := middleware.GetUserID(c)
			return h.projectService.GetProjects(c, userID, query)
		},
		http.StatusOK,
		&project.GetProjectsQuery{},
	)(c)
}

func (h *ProjectHandler) UpdateProject(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *project.UpdateProjectPayload) (*project.Project, error) {
			userID := middleware.GetUserID(c)
			return h.projectService.UpdateProject(c, userID, payload.ID, payload)
		},
		http.StatusOK,
		&project.UpdateProjectPayload{},
	)(c)
}

func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *project.DeleteProjectPayload) error {
			userID := middleware.GetUserID(c)
			return h.projectService.DeleteProject(c, userID, payload.ID)
		},
		http.StatusNoContent,
		&project.DeleteProjectPayload{},
	)(c)
}
