package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"taskflow/internal/errs"
	"taskflow/internal/server"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	server *server.Server
}

func NewAuthMiddleware(s *server.Server) *AuthMiddleware {
	return &AuthMiddleware{
		server: s,
	}
}

func (auth *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return auth.unauthorizedResponse(c, start, "missing authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return auth.unauthorizedResponse(c, start, "malformed authorization header")
		}

		claims, err := auth.server.JWT.Verify(parts[1])
		if err != nil {
			auth.server.Logger.Error().
				Err(err).
				Str("function", "RequireAuth").
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("could not verify token")
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		auth.server.Logger.Info().
			Str("function", "RequireAuth").
			Str("user_id", claims.UserID).
			Str("request_id", GetRequestID(c)).
			Dur("duration", time.Since(start)).
			Msg("user authenticated successfully")

		return next(c)
	}
}
func (auth *AuthMiddleware) unauthorizedResponse(c echo.Context, start time.Time, reason string) error {
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusUnauthorized)

	response := map[string]string{
		"code":     "UNAUTHORIZED",
		"message":  "Unauthorized",
		"override": "false",
		"status":   "401",
	}

	if err := json.NewEncoder(c.Response()).Encode(response); err != nil {
		auth.server.Logger.Error().
			Err(err).
			Str("function", "RequireAuth").
			Dur("duration", time.Since(start)).
			Msg("failed to write JSON response")
	} else {
		auth.server.Logger.Error().
			Str("function", "RequireAuth").
			Str("reason", reason).
			Dur("duration", time.Since(start)).
			Msg("could not get session claims from context")
	}

	return nil
}

func GetUserID(c echo.Context) string {
	if userID, ok := c.Get("user_id").(string); ok {
		return userID
	}
	return ""
}

func GetEmail(c echo.Context) string {
	if email, ok := c.Get("email").(string); ok {
		return email
	}
	return ""
}
