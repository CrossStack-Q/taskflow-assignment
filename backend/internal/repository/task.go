package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"taskflow/internal/model/task"
	"time"

	"taskflow/internal/server"

	"taskflow/internal/errs"
	"taskflow/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TaskRepository struct {
	server *server.Server
}

func NewTaskRepository(server *server.Server) *TaskRepository {
	return &TaskRepository{server: server}
}

func (r *TaskRepository) CreateTask(ctx context.Context, userID string, payload *task.CreateTaskPayload) (*task.Task, error) {

	stmt := `
		INSERT INTO
			tasks (
				user_id,
				title,
				description,
				priority,
				due_date,
				parent_task_id,
				project_id,
				metadata
			)
		VALUES
			(
				@user_id,
				@title,
				@description,
				@priority,
				@due_date,
				@parent_task_id,
				@project_id,
				@metadata
			)
		RETURNING
		*
	`
	priority := task.PriorityMedium
	if payload.Priority != nil {
		priority = *payload.Priority
	}

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":        userID,
		"title":          payload.Title,
		"description":    payload.Description,
		"priority":       priority,
		"due_date":       payload.DueDate,
		"parent_task_id": payload.ParentTaskID,
		"project_id":     payload.ProjectID,
		"metadata":       payload.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create task query for user_id=%s title=%s: %w", userID, payload.Title, err)
	}

	taskItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[task.Task])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:tasks for user_id=%s title=%s: %w", userID, payload.Title, err)
	}

	return &taskItem, nil
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, userID string, taskID uuid.UUID) (*task.PopulatedTask, error) {
	stmt := `
	SELECT
		t.*,
		CASE
			WHEN c.id IS NOT NULL THEN to_jsonb(camel (c))
			ELSE NULL
		END AS project,
		COALESCE(
			jsonb_agg(
				to_jsonb(camel (child))
				ORDER BY
					child.sort_order ASC,
					child.created_at ASC
			) FILTER (
				WHERE
					child.id IS NOT NULL
			),
			'[]'::JSONB
		) AS children,
		COALESCE(
			jsonb_agg(
				to_jsonb(camel (com))
				ORDER BY
					com.created_at ASC
			) FILTER (
				WHERE
					com.id IS NOT NULL
			),
			'[]'::JSONB
		) AS comments
	FROM
		tasks t
		LEFT JOIN task_projects c ON c.id=t.project_id
		AND c.user_id=@user_id
		LEFT JOIN tasks child ON child.parent_task_id=t.id
		AND child.user_id=@user_id
		LEFT JOIN task_comments com ON com.task_id=t.id
	WHERE
		t.id=@id
		AND t.user_id=@user_id
	GROUP BY
		t.id,
		c.id
`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      taskID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get task by id query for task_id=%s user_id=%s: %w", taskID.String(), userID, err)
	}

	taskItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[task.PopulatedTask])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:tasks for task_id=%s user_id=%s: %w", taskID.String(), userID, err)
	}

	return &taskItem, nil
}

func (r *TaskRepository) CheckTaskExists(ctx context.Context, userID string, taskID uuid.UUID) (*task.Task, error) {
	stmt := `
		SELECT
			*
		FROM
			tasks
		WHERE
			id=@id
			AND user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      taskID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check if task exists for task_id=%s user_id=%s: %w", taskID.String(), userID, err)
	}

	taskItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[task.Task])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:tasks for task_id=%s user_id=%s: %w", taskID.String(), userID, err)
	}

	return &taskItem, nil
}

func (r *TaskRepository) GetTasks(ctx context.Context, userID string, query *task.GetTasksQuery) (*model.PaginatedResponse[task.PopulatedTask], error) {
	stmt := `
	SELECT
		t.*,
		CASE
			WHEN c.id IS NOT NULL THEN to_jsonb(camel (c))
			ELSE NULL
		END AS project,
		COALESCE(
			jsonb_agg(
				to_jsonb(camel (child))
				ORDER BY
					child.sort_order ASC,
					child.created_at ASC
			) FILTER (
				WHERE
					child.id IS NOT NULL
			),
			'[]'::JSONB
		) AS children,
		COALESCE(
			jsonb_agg(
				to_jsonb(camel (com))
				ORDER BY
					com.created_at ASC
			) FILTER (
				WHERE
					com.id IS NOT NULL
			),
			'[]'::JSONB
		) AS comments
	FROM
		tasks t
		LEFT JOIN task_projects c ON c.id=t.project_id
		AND c.user_id=@user_id
		LEFT JOIN tasks child ON child.parent_task_id=t.id
		AND child.user_id=@user_id
		LEFT JOIN task_comments com ON com.task_id=t.id
`

	args := pgx.NamedArgs{
		"user_id": userID,
	}
	conditions := []string{"t.user_id = @user_id"}

	if query.Status != nil {
		conditions = append(conditions, "t.status = @status")
		args["status"] = *query.Status
	}

	if query.Priority != nil {
		conditions = append(conditions, "t.priority = @priority")
		args["priority"] = *query.Priority
	}

	if query.ProjectID != nil {
		conditions = append(conditions, "t.project_id = @project_id")
		args["project_id"] = *query.ProjectID
	}

	if query.ParentTaskID != nil {
		conditions = append(conditions, "t.parent_task_id = @parent_task_id")
		args["parent_task_id"] = *query.ParentTaskID
	} else {
		conditions = append(conditions, "t.parent_task_id IS NULL")
	}

	if query.DueFrom != nil {
		conditions = append(conditions, "t.due_date >= @due_from")
		args["due_from"] = *query.DueFrom
	}

	if query.DueTo != nil {
		conditions = append(conditions, "t.due_date <= @due_to")
		args["due_to"] = *query.DueTo
	}

	if query.Overdue != nil && *query.Overdue {
		conditions = append(conditions, "t.due_date < NOW() AND t.status != 'completed'")
	}

	if query.Completed != nil {
		if *query.Completed {
			conditions = append(conditions, "t.status = 'completed'")
		} else {
			conditions = append(conditions, "t.status != 'completed'")
		}
	}

	if query.Search != nil {
		conditions = append(conditions, "(t.title ILIKE @search OR t.description ILIKE @search)")
		args["search"] = "%" + *query.Search + "%"
	}

	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	countStmt := "SELECT COUNT(*) FROM tasks t"
	if len(conditions) > 0 {
		countStmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := r.server.DB.Pool.QueryRow(ctx, countStmt, args).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count for tasks user_id=%s: %w", userID, err)
	}

	stmt += " GROUP BY t.id, c.id"

	if query.Sort != nil {
		stmt += " ORDER BY t." + *query.Sort
		if query.Order != nil && *query.Order == "desc" {
			stmt += " DESC"
		} else {
			stmt += " ASC"
		}
	} else {
		stmt += " ORDER BY t.created_at DESC"
	}

	stmt += " LIMIT @limit OFFSET @offset"
	args["limit"] = *query.Limit
	args["offset"] = (*query.Page - 1) * (*query.Limit)

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get tasks query for user_id=%s: %w", userID, err)
	}

	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[task.PopulatedTask])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[task.PopulatedTask]{
				Data:       []task.PopulatedTask{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:tasks for user_id=%s: %w", userID, err)
	}

	return &model.PaginatedResponse[task.PopulatedTask]{
		Data:       tasks,
		Page:       *query.Page,
		Limit:      *query.Limit,
		Total:      total,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, userID string, payload *task.UpdateTaskPayload) (*task.Task, error) {
	stmt := "UPDATE tasks SET "
	args := pgx.NamedArgs{
		"task_id": payload.ID,
		"user_id": userID,
	}
	setClauses := []string{}

	if payload.Title != nil {
		setClauses = append(setClauses, "title = @title")
		args["title"] = *payload.Title
	}

	if payload.Description != nil {
		setClauses = append(setClauses, "description = @description")
		args["description"] = *payload.Description
	}

	if payload.Status != nil {
		setClauses = append(setClauses, "status = @status")
		args["status"] = *payload.Status

		if *payload.Status == task.StatusCompleted {
			setClauses = append(setClauses, "completed_at = @completed_at")
			args["completed_at"] = time.Now()
		} else if *payload.Status != task.StatusCompleted {
			setClauses = append(setClauses, "completed_at = NULL")
		}
	}

	if payload.Priority != nil {
		setClauses = append(setClauses, "priority = @priority")
		args["priority"] = *payload.Priority
	}

	if payload.DueDate != nil {
		setClauses = append(setClauses, "due_date = @due_date")
		args["due_date"] = *payload.DueDate
	}

	if payload.ParentTaskID != nil {
		setClauses = append(setClauses, "parent_task_id = @parent_task_id")
		args["parent_task_id"] = *payload.ParentTaskID
	}

	if payload.ProjectID != nil {
		setClauses = append(setClauses, "project_id = @project_id")
		args["project_id"] = *payload.ProjectID
	}

	if payload.Metadata != nil {
		setClauses = append(setClauses, "metadata = @metadata")
		args["metadata"] = payload.Metadata
	}

	if len(setClauses) == 0 {
		return nil, errs.NewBadRequestError("no fields to update", false, nil, nil, nil)
	}

	stmt += strings.Join(setClauses, ", ")
	stmt += " WHERE id = @task_id AND user_id = @user_id RETURNING *"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	updatedTask, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[task.Task])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:tasks: %w", err)
	}

	return &updatedTask, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, userID string, taskID uuid.UUID) error {
	stmt := `
		DELETE FROM tasks
		WHERE
			id=@task_id
			AND user_id=@user_id
	`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"task_id": taskID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if result.RowsAffected() == 0 {
		code := "TASK_NOT_FOUND"
		return errs.NewNotFoundError("task not found", false, &code)
	}

	return nil
}

func (r *TaskRepository) GetTaskStats(ctx context.Context, userID string) (*task.TaskStats, error) {
	stmt := `
		SELECT
			COUNT(*) AS total,
			COUNT(
				CASE
					WHEN status='draft' THEN 1
				END
			) AS draft,
			COUNT(
				CASE
					WHEN status='active' THEN 1
				END
			) AS active,
			COUNT(
				CASE
					WHEN status='completed' THEN 1
				END
			) AS completed,
			COUNT(
				CASE
					WHEN status='archived' THEN 1
				END
			) AS archived,
			COUNT(
				CASE
					WHEN due_date<NOW()
					AND status!='completed' THEN 1
				END
			) AS overdue
		FROM
			tasks
		WHERE
			user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	stats, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[task.TaskStats])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:tasks: %w", err)
	}

	return &stats, nil
}
