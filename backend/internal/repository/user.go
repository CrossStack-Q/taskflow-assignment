package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"taskflow/internal/errs"
	"taskflow/internal/model/user"
	"taskflow/internal/server"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	server *server.Server
}

func NewUserRepository(server *server.Server) *UserRepository {
	return &UserRepository{server: server}
}

func (r *UserRepository) CreateUser(ctx context.Context, payload *user.CreateUserPayload) (*user.User, error) {
	stmt := `
		INSERT INTO users (
			name,
			email,
			password
		)
		VALUES (
			@name,
			@email,
			@password
		)
		RETURNING *
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"name":     payload.Name,
		"email":    payload.Email,
		"password": payload.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create user query for email=%s: %w", payload.Email, err)
	}

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:users for email=%s: %w", payload.Email, err)
	}

	return &created, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	stmt := `
		SELECT *
		FROM   users
		WHERE  id = @id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get user by id query for user_id=%s: %w", userID, err)
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		if isNoRows(err) {
			code := "USER_NOT_FOUND"
			return nil, errs.NewNotFoundError("user not found", false, &code)
		}
		return nil, fmt.Errorf("failed to collect row from table:users for user_id=%s: %w", userID, err)
	}

	return &u, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	stmt := `
		SELECT *
		FROM   users
		WHERE  email = @email
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"email": email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get user by email query for email=%s: %w", email, err)
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		if isNoRows(err) {
			code := "USER_NOT_FOUND"
			return nil, errs.NewNotFoundError("user not found", false, &code)
		}
		return nil, fmt.Errorf("failed to collect row from table:users for email=%s: %w", email, err)
	}

	return &u, nil
}

func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	stmt := `
		SELECT EXISTS (
			SELECT 1
			FROM   users
			WHERE  email = @email
		)
	`

	var exists bool
	err := r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"email": email,
	}).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence for email=%s: %w", email, err)
	}

	return exists, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userID uuid.UUID, payload *user.UpdateUserPayload) (*user.User, error) {
	args := pgx.NamedArgs{"user_id": userID}
	setClauses := []string{}

	if payload.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *payload.Name
	}

	if payload.Email != nil {
		setClauses = append(setClauses, "email = @email")
		args["email"] = *payload.Email
	}

	if payload.Password != nil {
		setClauses = append(setClauses, "password = @password")
		args["password"] = *payload.Password
	}

	if len(setClauses) == 0 {
		return nil, errs.NewBadRequestError("no fields to update", false, nil, nil, nil)
	}

	stmt := "UPDATE users SET " +
		strings.Join(setClauses, ", ") +
		" WHERE id = @user_id RETURNING *"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update user query for user_id=%s: %w", userID, err)
	}

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		if isNoRows(err) {
			code := "USER_NOT_FOUND"
			return nil, errs.NewNotFoundError("user not found", false, &code)
		}
		return nil, fmt.Errorf("failed to collect row from table:users for user_id=%s: %w", userID, err)
	}

	return &updated, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	stmt := `
		DELETE FROM users
		WHERE  id = @user_id
	`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to execute delete user query for user_id=%s: %w", userID, err)
	}

	if result.RowsAffected() == 0 {
		code := "USER_NOT_FOUND"
		return errs.NewNotFoundError("user not found", false, &code)
	}

	return nil
}

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
