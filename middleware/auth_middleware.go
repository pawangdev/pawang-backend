package middleware

import (
	"pawang-backend/config"
	"pawang-backend/exception"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func Authenticated() func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(exception.ResponseError(false, "Unauthorized", fiber.StatusUnauthorized, nil))
		},
		SigningKey: []byte(config.GetEnv("JWT_SECRET_KEY")),
	})
}
