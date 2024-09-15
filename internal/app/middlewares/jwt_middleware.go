package middlewares

import (
	"github.com/labstack/echo/v4"
	"mamlaka/internal/pkg/tokens"
	"net/http"
	"strings"
)

// JWTMiddleware is a middleware for validating JWT access tokens.
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization Header")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization Header Format")
		}

		_, err := tokens.ValidateToken(tokenString, false) // false indicates it's an access token
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or Expired Token")
		}

		return next(c)
	}
}
