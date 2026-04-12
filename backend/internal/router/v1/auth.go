package v1

import (
	"taskflow/internal/handler"
	"taskflow/internal/middleware"

	"github.com/labstack/echo/v4"
)

func registerAuthRoutes(r *echo.Group, h *handler.AuthHandler) {
	auth := r.Group("/auth")

	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
}

func registerUserRoutes(r *echo.Group, h *handler.AuthHandler, auth *middleware.AuthMiddleware) {
	user := r.Group("/user")
	user.Use(auth.RequireAuth)

	user.GET("", h.GetUser)
	user.PATCH("", h.UpdateUser)
}
