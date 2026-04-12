package comment

import (
	"taskflow/internal/model"

	"github.com/google/uuid"
)

type Comment struct {
	model.Base
	TaskID  uuid.UUID `json:"taskId" db:"task_id"`
	UserID  string    `json:"userId" db:"user_id"`
	Content string    `json:"content" db:"content"`
}
