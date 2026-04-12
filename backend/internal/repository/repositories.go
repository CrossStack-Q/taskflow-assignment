package repository

import "taskflow/internal/server"

type Repositories struct {
	Task    *TaskRepository
	Comment *CommentRepository
	Project *ProjectRepository
	User    *UserRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Task:    NewTaskRepository(s),
		Comment: NewCommentRepository(s),
		Project: NewProjectRepository(s),
		User:    NewUserRepository(s),
	}
}
