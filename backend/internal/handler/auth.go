package handler

import (
	"net/http"

	"taskflow/internal/middleware"
	"taskflow/internal/model/user"
	"taskflow/internal/server"
	"taskflow/internal/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Handler
	authService *service.AuthService
}

func NewAuthHandler(s *server.Server, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		Handler:     NewHandler(s),
		authService: authService,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *user.RegisterPayload) (*user.AuthResponse, error) {
			u, token, err := h.authService.Register(c, payload)
			if err != nil {
				return nil, err
			}
			return &user.AuthResponse{User: u, Token: token}, nil
		},
		http.StatusCreated,
		&user.RegisterPayload{},
	)(c)
}

func (h *AuthHandler) Login(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *user.LoginPayload) (*user.AuthResponse, error) {
			u, token, err := h.authService.Login(c, payload)
			if err != nil {
				return nil, err
			}
			return &user.AuthResponse{User: u, Token: token}, nil
		},
		http.StatusOK,
		&user.LoginPayload{},
	)(c)
}

func (h *AuthHandler) GetUser(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, _ *user.GetUserPayload) (*user.User, error) {
			userID, err := uuid.Parse(middleware.GetUserID(c))
			if err != nil {
				return nil, err
			}
			return h.authService.Get(c, userID)
		},
		http.StatusOK,
		&user.GetUserPayload{},
	)(c)
}

func (h *AuthHandler) UpdateUser(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *user.UpdatePayload) (*user.User, error) {
			userID, err := uuid.Parse(middleware.GetUserID(c))
			if err != nil {
				return nil, err
			}
			return h.authService.Update(c, userID, payload)
		},
		http.StatusOK,
		&user.UpdatePayload{},
	)(c)
}
