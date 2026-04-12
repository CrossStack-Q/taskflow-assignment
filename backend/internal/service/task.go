package service

import (
	"taskflow/internal/errs"
	"taskflow/internal/middleware"
	"taskflow/internal/model"
	"taskflow/internal/model/task"
	"taskflow/internal/repository"
	"taskflow/internal/server"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TaskService struct {
	server      *server.Server
	taskRepo    *repository.TaskRepository
	projectRepo *repository.ProjectRepository
}

func NewTaskService(server *server.Server, taskRepo *repository.TaskRepository,
	projectRepo *repository.ProjectRepository,
) *TaskService {
	return &TaskService{
		server:      server,
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (s *TaskService) CreateTask(ctx echo.Context, userID string, payload *task.CreateTaskPayload) (*task.Task, error) {
	logger := middleware.GetLogger(ctx)

	if payload.ParentTaskID != nil {
		parentTask, err := s.taskRepo.CheckTaskExists(ctx.Request().Context(), userID, *payload.ParentTaskID)
		if err != nil {
			logger.Error().Err(err).Msg("parent task validation failed")
			return nil, err
		}

		if !parentTask.CanHaveChildren() {
			err := errs.NewBadRequestError("Parent task cannot have children (subtasks can't have subtasks)", false, nil, nil, nil)
			logger.Warn().Msg("parent task cannot have children")
			return nil, err
		}
	}

	if payload.ProjectID != nil {
		_, err := s.projectRepo.GetProjectByID(ctx.Request().Context(), userID, *payload.ProjectID)
		if err != nil {
			logger.Error().Err(err).Msg("project validation failed")
			return nil, err
		}
	}

	taskItem, err := s.taskRepo.CreateTask(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create task")
		return nil, err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "task_created").
		Str("task_id", taskItem.ID.String()).
		Str("title", taskItem.Title).
		Str("project_id", func() string {
			if taskItem.ProjectID != nil {
				return taskItem.ProjectID.String()
			}
			return ""
		}()).
		Str("priority", string(taskItem.Priority)).
		Msg("Task created successfully")

	return taskItem, nil
}

func (s *TaskService) GetTaskByID(ctx echo.Context, userID string, taskID uuid.UUID) (*task.PopulatedTask, error) {
	logger := middleware.GetLogger(ctx)

	taskItem, err := s.taskRepo.GetTaskByID(ctx.Request().Context(), userID, taskID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch task by ID")
		return nil, err
	}

	return taskItem, nil
}

func (s *TaskService) GetTasks(ctx echo.Context, userID string, query *task.GetTasksQuery) (*model.PaginatedResponse[task.PopulatedTask], error) {
	logger := middleware.GetLogger(ctx)

	result, err := s.taskRepo.GetTasks(ctx.Request().Context(), userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch tasks")
		return nil, err
	}

	return result, nil
}

func (s *TaskService) UpdateTask(ctx echo.Context, userID string, payload *task.UpdateTaskPayload) (*task.Task, error) {
	logger := middleware.GetLogger(ctx)

	if payload.ParentTaskID != nil {
		parentTask, err := s.taskRepo.CheckTaskExists(ctx.Request().Context(), userID, *payload.ParentTaskID)
		if err != nil {
			logger.Error().Err(err).Msg("parent task validation failed")
			return nil, err
		}

		if parentTask.ID == payload.ID {
			err := errs.NewBadRequestError("Task cannot be its own parent", false, nil, nil, nil)
			logger.Warn().Msg("task cannot be its own parent")
			return nil, err
		}

		if !parentTask.CanHaveChildren() {
			err := errs.NewBadRequestError("Parent task cannot have children (subtasks can't have subtasks)", false, nil, nil, nil)
			logger.Warn().Msg("parent task cannot have children")
			return nil, err
		}

		logger.Debug().Msg("parent task validation passed")
	}

	if payload.ProjectID != nil {
		_, err := s.projectRepo.GetProjectByID(ctx.Request().Context(), userID, *payload.ProjectID)
		if err != nil {
			logger.Error().Err(err).Msg("project validation failed")
			return nil, err
		}

		logger.Debug().Msg("project validation passed")
	}

	updatedTask, err := s.taskRepo.UpdateTask(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update task")
		return nil, err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "task_updated").
		Str("task_id", updatedTask.ID.String()).
		Str("title", updatedTask.Title).
		Str("project_id", func() string {
			if updatedTask.ProjectID != nil {
				return updatedTask.ProjectID.String()
			}
			return ""
		}()).
		Str("priority", string(updatedTask.Priority)).
		Str("status", string(updatedTask.Status)).
		Msg("Task updated successfully")

	return updatedTask, nil
}

func (s *TaskService) DeleteTask(ctx echo.Context, userID string, taskID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	err := s.taskRepo.DeleteTask(ctx.Request().Context(), userID, taskID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete task")
		return err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "task_deleted").
		Str("task_id", taskID.String()).
		Msg("Task deleted successfully")

	return nil
}

func (s *TaskService) GetTaskStats(ctx echo.Context, userID string) (*task.TaskStats, error) {
	logger := middleware.GetLogger(ctx)

	stats, err := s.taskRepo.GetTaskStats(ctx.Request().Context(), userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch task statistics")
		return nil, err
	}

	return stats, nil
}
