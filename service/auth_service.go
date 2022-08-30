package service

import (
	"pawang-backend/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	GenerateToken(userID int) (string, error)
	CurrentLoggedUserID(c *fiber.Ctx) int
}

type authService struct {
}

func NewAuthService() *authService {
	return &authService{}
}

func (service *authService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET_KEY")))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (service *authService) CurrentLoggedUserID(c *fiber.Ctx) int {
	user := c.Locals("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	userID := claim["user_id"].(float64)
	return int(userID)
}
