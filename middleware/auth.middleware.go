package middleware

import (
	"pawang-backend/config"

	"github.com/labstack/echo/v4/middleware"
)

var IsAuthenticated = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(config.GetEnv("JWT_TOKEN_SECRET")),
})
