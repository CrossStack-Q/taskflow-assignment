package service

import (
	"taskflow/internal/repository"
	"taskflow/internal/server"
)

type Services struct {
	Task    *TaskService
	Comment *CommentService
	Project *ProjectService
	Auth    *AuthService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {

	return &Services{
		Project: NewProjectService(s, repos.Project),
		Comment: NewCommentService(s, repos.Comment, repos.Task),
		Task:    NewTaskService(s, repos.Task, repos.Project),
		Auth:    NewAuthService(s, repos.User),
	}, nil
}
