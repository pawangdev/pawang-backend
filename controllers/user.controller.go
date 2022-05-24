package controllers

import (
	"net/http"
	"pawang-backend/config"
	"pawang-backend/helpers"
	"pawang-backend/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type inputRegister struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Phone    string `json:"phone" form:"phone"`
	Gender   string `json:"gender" form:"gender"`
}

type inputLogin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func Register(c echo.Context) error {
	db := config.ConnectDatabase()
	var input inputRegister
	user := new(models.User)

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	hashed, err := helpers.HashPassword(input.Password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Password = hashed
	user.Phone = input.Phone
	user.Gender = input.Gender

	if err := db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, models.Response{Success: true, Data: nil, Message: "register successfully, please login your account"})
}

func Login(c echo.Context) error {
	db := config.ConnectDatabase()
	var input inputLogin
	user := new(models.User)

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	result := db.Find(&user, "email = ?", input.Email)

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Success: false, Data: nil, Message: "email tidak dapat ditemukan"})
	}

	check, _ := helpers.CompareHashPassword(input.Password, user.Password)

	if !check {
		return c.JSON(http.StatusUnauthorized, models.Response{Success: false, Data: nil, Message: "password salah"})
	}

	claims := jwt.MapClaims{}
	claims["user_id"] = user.ID
	claims["name"] = user.Name
	claims["email"] = user.Email

	token, err := helpers.GenerateToken(claims)
	if err != nil {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusAccepted, models.Response{Success: true, Data: token, Message: "login success"})
}

func Profile(c echo.Context) error {
	db := config.ConnectDatabase()
	user := new(models.User)

	userID := helpers.GetLoginUserID(c)

	result := db.Find(&user, "id = ?", userID)

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusUnauthorized, models.Response{Success: false, Data: nil})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Data: user})
}
