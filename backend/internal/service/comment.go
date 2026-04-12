package service

import (
	"taskflow/internal/middleware"
	"taskflow/internal/model/comment"
	"taskflow/internal/repository"
	"taskflow/internal/server"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CommentService struct {
	server      *server.Server
	commentRepo *repository.CommentRepository
	taskRepo    *repository.TaskRepository
}

func NewCommentService(server *server.Server, commentRepo *repository.CommentRepository, taskRepo *repository.TaskRepository) *CommentService {
	return &CommentService{
		server:      server,
		commentRepo: commentRepo,
		taskRepo:    taskRepo,
	}
}

func (s *CommentService) AddComment(ctx echo.Context, userID string, taskID uuid.UUID,
	payload *comment.AddCommentPayload,
) (*comment.Comment, error) {
	logger := middleware.GetLogger(ctx)

	_, err := s.taskRepo.CheckTaskExists(ctx.Request().Context(), userID, taskID)
	if err != nil {
		logger.Error().Err(err).Msg("task validation failed")
		return nil, err
	}

	commentItem, err := s.commentRepo.AddComment(ctx.Request().Context(), userID, taskID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to add comment")
		return nil, err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "comment_added").
		Str("comment_id", commentItem.ID.String()).
		Str("task_id", taskID.String()).
		Msg("Comment added successfully")

	return commentItem, nil
}

func (s *CommentService) GetCommentsByTaskID(ctx echo.Context, userID string, taskID uuid.UUID) ([]comment.Comment, error) {
	logger := middleware.GetLogger(ctx)

	_, err := s.taskRepo.CheckTaskExists(ctx.Request().Context(), userID, taskID)
	if err != nil {
		logger.Error().Err(err).Msg("task validation failed")
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByTaskID(ctx.Request().Context(), userID, taskID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch comments by task ID")
		return nil, err
	}

	return comments, nil
}

func (s *CommentService) UpdateComment(ctx echo.Context, userID string, commentID uuid.UUID, content string) (*comment.Comment, error) {
	logger := middleware.GetLogger(ctx)

	_, err := s.commentRepo.GetCommentByID(ctx.Request().Context(), userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("comment validation failed")
		return nil, err
	}

	commentItem, err := s.commentRepo.UpdateComment(ctx.Request().Context(), userID, commentID, content)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update comment")
		return nil, err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "comment_updated").
		Str("comment_id", commentItem.ID.String()).
		Msg("Comment updated successfully")

	return commentItem, nil
}

func (s *CommentService) DeleteComment(ctx echo.Context, userID string, commentID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	_, err := s.commentRepo.GetCommentByID(ctx.Request().Context(), userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("comment validation failed")
		return err
	}

	err = s.commentRepo.DeleteComment(ctx.Request().Context(), userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete comment")
		return err
	}
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "comment_deleted").
		Str("comment_id", commentID.String()).
		Msg("Comment deleted successfully")

	return nil
}
