package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"taskflow/internal/model/project"

	"taskflow/internal/model"
	"taskflow/internal/server"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProjectRepository struct {
	server *server.Server
}

func NewProjectRepository(server *server.Server) *ProjectRepository {
	return &ProjectRepository{server: server}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, userID string,
	payload *project.CreateProjectPayload,
) (*project.Project, error) {
	stmt := `
		INSERT INTO
			task_projects (
				user_id,
				name,
				color,
				description
			)
		VALUES
			(
				@user_id,
				@name,
				@color,
				@description
			)
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":     userID,
		"name":        payload.Name,
		"color":       payload.Color,
		"description": payload.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create project query for user_id=%s name=%s: %w", userID, payload.Name, err)
	}

	projectItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[project.Project])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:task_projects for user_id=%s name=%s: %w", userID, payload.Name, err)
	}

	return &projectItem, nil
}

func (r *ProjectRepository) GetProjectByID(ctx context.Context, userID string, projectID uuid.UUID) (*project.Project, error) {
	stmt := `
		SELECT
			*
		FROM
			task_projects
		WHERE
			id=@id
			AND user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      projectID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get project by id query for project_id=%s user_id=%s: %w", projectID.String(), userID, err)
	}

	projectItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[project.Project])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:task_projects for project_id=%s user_id=%s: %w", projectID.String(), userID, err)
	}

	return &projectItem, nil
}

func (r *ProjectRepository) GetProjects(ctx context.Context, userID string,
	query *project.GetProjectsQuery,
) (*model.PaginatedResponse[project.Project], error) {
	stmt := `
		SELECT
			*
		FROM
			task_projects
		WHERE
			user_id=@user_id
	`

	args := pgx.NamedArgs{
		"user_id": userID,
	}

	if query.Search != nil {
		stmt += ` AND name ILIKE '%' || @search || '%'`
		args["search"] = *query.Search
	}

	sortColumn := "name"
	if query.Sort != nil {
		sortColumn = *query.Sort
	}
	sortOrder := "asc"
	if query.Order != nil {
		sortOrder = *query.Order
	}
	stmt += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)

	stmt += ` LIMIT @limit OFFSET @offset`
	args["limit"] = *query.Limit
	args["offset"] = (*query.Page - 1) * (*query.Limit)

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get projects query for user_id=%s: %w", userID, err)
	}

	projects, err := pgx.CollectRows(rows, pgx.RowToStructByName[project.Project])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[project.Project]{
				Data:       []project.Project{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:task_projects for user_id=%s: %w", userID, err)
	}

	countStmt := `
		SELECT
			COUNT(*)
		FROM
			task_projects
		WHERE
			user_id=@user_id
	`

	countArgs := pgx.NamedArgs{
		"user_id": userID,
	}

	if query.Search != nil {
		countStmt += ` AND name ILIKE '%' || @search || '%'`
		countArgs["search"] = *query.Search
	}

	var total int
	err = r.server.DB.Pool.QueryRow(ctx, countStmt, countArgs).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count of projects for user_id=%s: %w", userID, err)
	}

	return &model.PaginatedResponse[project.Project]{
		Data:       projects,
		Page:       *query.Page,
		Limit:      *query.Limit,
		Total:      total,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, userID string,
	projectID uuid.UUID, payload *project.UpdateProjectPayload,
) (*project.Project, error) {
	stmt := `UPDATE task_projects SET `
	args := pgx.NamedArgs{
		"id":      projectID,
		"user_id": userID,
	}
	setClauses := []string{}

	if payload.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *payload.Name
	}
	if payload.Color != nil {
		setClauses = append(setClauses, "color = @color")
		args["color"] = *payload.Color
	}
	if payload.Description != nil {
		setClauses = append(setClauses, "description = @description")
		args["description"] = *payload.Description
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	stmt += strings.Join(setClauses, ", ")
	stmt += ` WHERE id = @id AND user_id = @user_id RETURNING *`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update project query for project_id=%s user_id=%s: %w", projectID.String(), userID, err)
	}

	projectItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[project.Project])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:task_projects for project_id=%s user_id=%s: %w", projectID.String(), userID, err)
	}

	return &projectItem, nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, userID string, projectID uuid.UUID) error {
	result, err := r.server.DB.Pool.Exec(ctx, `
		DELETE FROM task_projects
		WHERE id = @id AND user_id = @user_id
	`, pgx.NamedArgs{
		"id":      projectID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}
