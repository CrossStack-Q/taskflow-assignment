package service

import (
	"taskflow/internal/errs"
	"taskflow/internal/middleware"
	"taskflow/internal/model/user"
	"taskflow/internal/repository"
	"taskflow/internal/server"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	server   *server.Server
	userRepo *repository.UserRepository
}

func NewAuthService(server *server.Server, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		server:   server,
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(ctx echo.Context, payload *user.RegisterPayload) (*user.User, string, error) {
	logger := middleware.GetLogger(ctx)

	exists, err := s.userRepo.CheckEmailExists(ctx.Request().Context(), payload.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to check email existence")
		return nil, "", err
	}
	if exists {
		code := "EMAIL_TAKEN"
		return nil, "", errs.NewConflictError("email is already in use", false, &code)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		logger.Error().Err(err).Msg("failed to hash password")
		return nil, "", err
	}

	created, err := s.userRepo.CreateUser(ctx.Request().Context(), &user.CreateUserPayload{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: string(hashed),
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create user")
		return nil, "", err
	}

	token, err := s.server.JWT.Sign(created.ID, created.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to sign JWT after registration")
		return nil, "", err
	}

	logger.Info().
		Str("event", "user_registered").
		Str("user_id", created.ID.String()).
		Str("email", created.Email).
		Msg("User registered successfully")

	return created, token, nil
}

func (s *AuthService) Login(ctx echo.Context, payload *user.LoginPayload) (*user.User, string, error) {
	logger := middleware.GetLogger(ctx)

	u, err := s.userRepo.GetUserByEmail(ctx.Request().Context(), payload.Email)
	if err != nil {
		logger.Warn().Err(err).Str("email", payload.Email).Msg("login failed: user lookup")
		return nil, "", errs.NewUnauthorizedError("invalid email or password", false)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(payload.Password)); err != nil {
		logger.Warn().Str("user_id", u.ID.String()).Msg("login failed: wrong password")

		return nil, "", errs.NewUnauthorizedError("invalid email or password", false)
	}

	token, err := s.server.JWT.Sign(u.ID, u.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to sign JWT after login")
		return nil, "", err
	}

	logger.Info().
		Str("event", "user_logged_in").
		Str("user_id", u.ID.String()).
		Str("email", u.Email).
		Msg("User logged in successfully")

	return u, token, nil
}

func (s *AuthService) Get(ctx echo.Context, userID uuid.UUID) (*user.User, error) {
	logger := middleware.GetLogger(ctx)

	u, err := s.userRepo.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID.String()).Msg("failed to fetch current user")
		return nil, err
	}

	return u, nil
}

func (s *AuthService) Update(ctx echo.Context, userID uuid.UUID, payload *user.UpdatePayload) (*user.User, error) {
	logger := middleware.GetLogger(ctx)

	if payload.Email != nil {
		exists, err := s.userRepo.CheckEmailExists(ctx.Request().Context(), *payload.Email)
		if err != nil {
			logger.Error().Err(err).Msg("failed to check email existence during update")
			return nil, err
		}
		if exists {
			code := "EMAIL_TAKEN"
			return nil, errs.NewConflictError("email is already in use", false, &code)
		}
	}

	repoPayload := &user.UpdateUserPayload{
		Name:  payload.Name,
		Email: payload.Email,
	}

	if payload.NewPassword != nil {
		current, err := s.userRepo.GetUserByID(ctx.Request().Context(), userID)
		if err != nil {
			logger.Error().Err(err).Msg("failed to fetch user for password change")
			return nil, err
		}

		if payload.CurrentPassword == nil {
			code := "CURRENT_PASSWORD_REQUIRED"
			return nil, errs.NewBadRequestError("current_password is required to set a new password", false, &code, nil, nil)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(current.Password), []byte(*payload.CurrentPassword)); err != nil {
			logger.Warn().Str("user_id", userID.String()).Msg("password change failed: wrong current password")
			code := "INVALID_CURRENT_PASSWORD"
			return nil, errs.NewBadRequestError("current password is incorrect", false, &code, nil, nil)
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(*payload.NewPassword), 12)
		if err != nil {
			logger.Error().Err(err).Msg("failed to hash new password")
			return nil, err
		}
		h := string(hashed)
		repoPayload.Password = &h
	}

	updated, err := s.userRepo.UpdateUser(ctx.Request().Context(), userID, repoPayload)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID.String()).Msg("failed to update user")
		return nil, err
	}

	logger.Info().
		Str("event", "user_updated").
		Str("user_id", updated.ID.String()).
		Msg("User profile updated successfully")

	return updated, nil
}
