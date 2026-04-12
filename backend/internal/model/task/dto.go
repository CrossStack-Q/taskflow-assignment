package task

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateTaskPayload struct {
	Title        string     `json:"title" validate:"required,min=1,max=255"`
	Description  *string    `json:"description" validate:"omitempty,max=1000"`
	Priority     *Priority  `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate      *time.Time `json:"dueDate"`
	ParentTaskID *uuid.UUID `json:"parentTaskId" validate:"omitempty,uuid"`
	ProjectID    *uuid.UUID `json:"projectId" validate:"omitempty,uuid"`
	Metadata     *Metadata  `json:"metadata"`
}

func (p *CreateTaskPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type UpdateTaskPayload struct {
	ID           uuid.UUID  `param:"id" validate:"required,uuid"`
	Title        *string    `json:"title" validate:"omitempty,min=1,max=255"`
	Description  *string    `json:"description" validate:"omitempty,max=1000"`
	Status       *Status    `json:"status" validate:"omitempty,oneof=draft active completed archived"`
	Priority     *Priority  `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate      *time.Time `json:"dueDate"`
	ParentTaskID *uuid.UUID `json:"parentTaskId" validate:"omitempty,uuid"`
	ProjectID    *uuid.UUID `json:"projectId" validate:"omitempty,uuid"`
	Metadata     *Metadata  `json:"metadata"`
}

func (p *UpdateTaskPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type GetTasksQuery struct {
	Page         *int       `query:"page" validate:"omitempty,min=1"`
	Limit        *int       `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort         *string    `query:"sort" validate:"omitempty,oneof=created_at updated_at title priority due_date status"`
	Order        *string    `query:"order" validate:"omitempty,oneof=asc desc"`
	Search       *string    `query:"search" validate:"omitempty,min=1"`
	Status       *Status    `query:"status" validate:"omitempty,oneof=draft active completed archived"`
	Priority     *Priority  `query:"priority" validate:"omitempty,oneof=low medium high"`
	ProjectID    *uuid.UUID `query:"projectId" validate:"omitempty,uuid"`
	ParentTaskID *uuid.UUID `query:"parentTaskId" validate:"omitempty,uuid"`
	DueFrom      *time.Time `query:"dueFrom"`
	DueTo        *time.Time `query:"dueTo"`
	Overdue      *bool      `query:"overdue"`
	Completed    *bool      `query:"completed"`
}

func (q *GetTasksQuery) Validate() error {
	validate := validator.New()

	if err := validate.Struct(q); err != nil {
		return err
	}

	if q.Page == nil {
		defaultPage := 1
		q.Page = &defaultPage
	}
	if q.Limit == nil {
		defaultLimit := 20
		q.Limit = &defaultLimit
	}
	if q.Sort == nil {
		defaultSort := "created_at"
		q.Sort = &defaultSort
	}
	if q.Order == nil {
		defaultOrder := "desc"
		q.Order = &defaultOrder
	}

	return nil
}

type GetTaskByIDPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *GetTaskByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type DeleteTaskPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *DeleteTaskPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type GetTaskStatsPayload struct{}

func (p *GetTaskStatsPayload) Validate() error {
	return nil
}

type UploadTaskAttachmentPayload struct {
	TaskID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *UploadTaskAttachmentPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type DeleteTaskAttachmentPayload struct {
	TaskID       uuid.UUID `param:"id" validate:"required,uuid"`
	AttachmentID uuid.UUID `param:"attachmentId" validate:"required,uuid"`
}

func (p *DeleteTaskAttachmentPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type GetAttachmentPresignedURLPayload struct {
	TaskID       uuid.UUID `param:"id" validate:"required,uuid"`
	AttachmentID uuid.UUID `param:"attachmentId" validate:"required,uuid"`
}

func (p *GetAttachmentPresignedURLPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
