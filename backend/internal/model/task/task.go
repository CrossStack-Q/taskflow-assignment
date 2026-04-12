package task

import (
	"time"

	"taskflow/internal/model"
	"taskflow/internal/model/comment"
	"taskflow/internal/model/project"

	"github.com/google/uuid"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
	StatusArchived  Status = "archived"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	model.Base
	UserID       string     `json:"userId" db:"user_id"`
	Title        string     `json:"title" db:"title"`
	Description  *string    `json:"description" db:"description"`
	Status       Status     `json:"status" db:"status"`
	Priority     Priority   `json:"priority" db:"priority"`
	DueDate      *time.Time `json:"dueDate" db:"due_date"`
	CompletedAt  *time.Time `json:"completedAt" db:"completed_at"`
	ParentTaskID *uuid.UUID `json:"parentTaskId" db:"parent_task_id"`
	ProjectID    *uuid.UUID `json:"projectId" db:"project_id"`
	Metadata     *Metadata  `json:"metadata" db:"metadata"`
	SortOrder    int        `json:"sortOrder" db:"sort_order"`
}

type Metadata struct {
	Tags       []string `json:"tags"`
	Reminder   *string  `json:"reminder"`
	Color      *string  `json:"color"`
	Difficulty *int     `json:"difficulty"`
}

type PopulatedTask struct {
	Task
	Project  *project.Project  `json:"project" db:"project"`
	Children []Task            `json:"children" db:"children"`
	Comments []comment.Comment `json:"comments" db:"comments"`
}

type TaskStats struct {
	Total     int `json:"total"`
	Draft     int `json:"draft"`
	Active    int `json:"active"`
	Completed int `json:"completed"`
	Archived  int `json:"archived"`
	Overdue   int `json:"overdue"`
}

type UserWeeklyStats struct {
	UserID         string `json:"userId" db:"user_id"`
	CreatedCount   int    `json:"createdCount" db:"created_count"`
	CompletedCount int    `json:"completedCount" db:"completed_count"`
	ActiveCount    int    `json:"activeCount" db:"active_count"`
	OverdueCount   int    `json:"overdueCount" db:"overdue_count"`
}

func (t *Task) IsOverdue() bool {
	return t.DueDate != nil && t.DueDate.Before(time.Now()) && t.Status != StatusCompleted
}

func (t *Task) CanHaveChildren() bool {
	return t.ParentTaskID == nil
}
