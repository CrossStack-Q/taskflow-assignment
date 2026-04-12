package v1

import (
	"taskflow/internal/handler"
	"taskflow/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterV1Routes(router *echo.Group, handlers *handler.Handlers, middleware *middleware.Middlewares) {

	registerAuthRoutes(router, handlers.Auth)
	registerUserRoutes(router, handlers.Auth, middleware.Auth)
	registerTaskRoutes(router, handlers.Task, handlers.Comment, middleware.Auth)

	registerProjectRoutes(router, handlers.Project, middleware.Auth)

	registerCommentRoutes(router, handlers.Comment, middleware.Auth)
}
