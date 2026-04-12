package service

import (
	"taskflow/internal/middleware"
	"taskflow/internal/model"
	"taskflow/internal/model/project"
	"taskflow/internal/repository"
	"taskflow/internal/server"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProjectService struct {
	server      *server.Server
	ProjectRepo *repository.ProjectRepository
}

func NewProjectService(server *server.Server, ProjectRepo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{
		server:      server,
		ProjectRepo: ProjectRepo,
	}
}

func (s *ProjectService) CreateProject(ctx echo.Context, userID string,
	payload *project.CreateProjectPayload,
) (*project.Project, error) {
	logger := middleware.GetLogger(ctx)

	ProjectItem, err := s.ProjectRepo.CreateProject(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create Project")
		return nil, err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "Project_created").
		Str("Project_id", ProjectItem.ID.String()).
		Str("name", ProjectItem.Name).
		Str("color", ProjectItem.Color).
		Msg("Project created successfully")

	return ProjectItem, nil
}

func (s *ProjectService) GetProjects(ctx echo.Context, userID string,
	query *project.GetProjectsQuery,
) (*model.PaginatedResponse[project.Project], error) {
	logger := middleware.GetLogger(ctx)

	projects, err := s.ProjectRepo.GetProjects(ctx.Request().Context(), userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch projects")
		return nil, err
	}

	return projects, nil
}

func (s *ProjectService) GetProjectByID(ctx echo.Context, userID string, ProjectID uuid.UUID) (*project.Project, error) {
	logger := middleware.GetLogger(ctx)

	ProjectItem, err := s.ProjectRepo.GetProjectByID(ctx.Request().Context(), userID, ProjectID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch Project by ID")
		return nil, err
	}

	return ProjectItem, nil
}

func (s *ProjectService) UpdateProject(ctx echo.Context, userID string, ProjectID uuid.UUID,
	payload *project.UpdateProjectPayload,
) (*project.Project, error) {
	logger := middleware.GetLogger(ctx)

	ProjectItem, err := s.ProjectRepo.UpdateProject(ctx.Request().Context(), userID, ProjectID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update Project")
		return nil, err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "Project_updated").
		Str("Project_id", ProjectItem.ID.String()).
		Str("name", ProjectItem.Name).
		Msg("Project updated successfully")

	return ProjectItem, nil
}

func (s *ProjectService) DeleteProject(ctx echo.Context, userID string, ProjectID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	err := s.ProjectRepo.DeleteProject(ctx.Request().Context(), userID, ProjectID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete Project")
		return err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "Project_deleted").
		Str("Project_id", ProjectID.String()).
		Msg("Project deleted successfully")

	return nil
}
