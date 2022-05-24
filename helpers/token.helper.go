package helpers

import (
	"fmt"
	"pawang-backend/config"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.GetEnv("JWT_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}

func GetLoginUserID(c echo.Context) uint {
	userLogin := c.Get("user").(*jwt.Token)
	claims := userLogin.Claims.(jwt.MapClaims)
	uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
	if err != nil {
		return 0
	}

	return uint(uid)
}
