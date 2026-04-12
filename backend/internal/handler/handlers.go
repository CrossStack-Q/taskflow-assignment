package handler

import (
	"taskflow/internal/server"
	"taskflow/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	Task    *TaskHandler
	Auth    *AuthHandler
	Comment *CommentHandler
	Project *ProjectHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		Task:    NewTaskHandler(s, services.Task),
		Project: NewProjectHandler(s, services.Project),
		Comment: NewCommentHandler(s, services.Comment),
		Auth:    NewAuthHandler(s, services.Auth),
	}
}
