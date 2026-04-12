package handler

import (
	"net/http"

	"taskflow/internal/middleware"
	"taskflow/internal/model/comment"
	"taskflow/internal/server"
	"taskflow/internal/service"

	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	Handler
	commentService *service.CommentService
}

func NewCommentHandler(s *server.Server, commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		Handler:        NewHandler(s),
		commentService: commentService,
	}
}

func (h *CommentHandler) AddComment(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *comment.AddCommentPayload) (*comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.commentService.AddComment(c, userID, payload.TaskID, payload)
		},
		http.StatusCreated,
		&comment.AddCommentPayload{},
	)(c)
}

func (h *CommentHandler) GetCommentsByTaskID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *comment.GetCommentsByTaskIDPayload) ([]comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.commentService.GetCommentsByTaskID(c, userID, payload.TaskID)
		},
		http.StatusOK,
		&comment.GetCommentsByTaskIDPayload{},
	)(c)
}

func (h *CommentHandler) UpdateComment(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *comment.UpdateCommentPayload) (*comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.commentService.UpdateComment(c, userID, payload.ID, payload.Content)
		},
		http.StatusOK,
		&comment.UpdateCommentPayload{},
	)(c)
}

func (h *CommentHandler) DeleteComment(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *comment.DeleteCommentPayload) error {
			userID := middleware.GetUserID(c)
			return h.commentService.DeleteComment(c, userID, payload.ID)
		},
		http.StatusNoContent,
		&comment.DeleteCommentPayload{},
	)(c)
}
