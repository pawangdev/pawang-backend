package controllers

import (
	"fmt"
	"net/http"
	"os"
	"pawang-backend/config"
	"pawang-backend/helpers"
	"pawang-backend/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	echogothic "github.com/nabowler/echo-gothic"
)

type inputRegister struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
	Phone    string `json:"phone" form:"phone" validate:"required"`
	Gender   string `json:"gender" form:"gender" validate:"required"`
}

type inputLogin struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type inputChangePassword struct {
	PasswordNow             string `json:"password_now" form:"password_now" validate:"required"`
	PasswordNew             string `json:"password_new" form:"password_new" validate:"required"`
	PasswordNewConfirmation string `json:"password_new_confirmation" form:"password_new_confirmation" validate:"required"`
}

type inputChangeProfile struct {
	Name string `json:"name" form:"name" validate:"required"`
}

func Register(c echo.Context) error {
	db := config.ConnectDatabase()
	var input inputRegister
	user := new(models.User)
	wallet := new(models.Wallet)

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

	wallet.Name = "Dompet"
	wallet.Balance = 0
	wallet.UserID = user.ID

	if err := db.Create(&wallet).Error; err != nil {
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

	return c.JSON(http.StatusAccepted, models.Response{Success: true, Data: map[string]string{"token": token}, Message: "login success"})
}

func LoginWithGoogle(c echo.Context) error {
	return echogothic.BeginAuthHandler(c)
}

func LoginWithGoogleCallback(c echo.Context) error {
	db := config.ConnectDatabase()
	userLogin := models.User{}

	user, err := echogothic.CompleteUserAuth(c)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Response{Success: false, Message: err.Error()})
	}

	result := db.Find(&userLogin, "email = ?", user.Email)

	if result.RowsAffected == 0 {
		hashed, _ := helpers.HashPassword(time.Now().String())
		userLogin.Name = user.Email
		userLogin.Email = user.Email
		userLogin.ImageProfile = user.AvatarURL
		userLogin.Password = hashed

		if err := db.Save(&userLogin).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
		}
	}

	claims := jwt.MapClaims{}
	claims["user_id"] = userLogin.ID
	claims["name"] = userLogin.Name
	claims["email"] = userLogin.Email

	token, err := helpers.GenerateToken(claims)
	if err != nil {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusAccepted, models.Response{Success: true, Data: map[string]string{"token": token}, Message: "login success"})
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

func ChangePassword(c echo.Context) error {
	db := config.ConnectDatabase()
	user := new(models.User)
	var input inputChangePassword

	results := db.Find(&user, "id = ?", helpers.GetLoginUserID(c))

	if results.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Success: false, Data: nil, Message: "Data Tidak Ditemukan"})
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil, Message: "Input Tidak Sesuai"})
	}

	checkOldPassword, _ := helpers.CompareHashPassword(input.PasswordNow, user.Password)

	if !checkOldPassword {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil, Message: "Password Lama Salah"})
	}

	if input.PasswordNew != input.PasswordNewConfirmation {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil, Message: "Password Baru Tidak Sesuai Dengan Password Konfirmasi"})
	}

	hashNewPassword, err := helpers.HashPassword(input.PasswordNew)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	user.Password = hashNewPassword

	if err := db.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Data: nil, Message: "Success Update Password"})
}

func UpdateProfile(c echo.Context) error {
	db := config.ConnectDatabase()
	user := new(models.User)
	var input inputChangeProfile
	results := db.Find(&user, "id = ?", helpers.GetLoginUserID(c))

	if results.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Success: false, Data: nil, Message: "Data Tidak Ditemukan"})
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil, Message: "Input Tidak Sesuai"})
	}

	if len(c.Request().MultipartForm.File) != 0 {
		// Get Form Image
		file, err := c.FormFile("image")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
		}

		// Setting Directory
		filePath := fmt.Sprintf("public/users/%v/profile_image/%v-%v", helpers.GetLoginUserID(c), time.Now().Unix(), file.Filename)
		fileSrc := fmt.Sprintf("profile_image/%v/profile_image/%v-%v", helpers.GetLoginUserID(c), time.Now().Unix(), file.Filename)
		dirPath := fmt.Sprintf("public/users/%v/profile_image/", helpers.GetLoginUserID(c))

		// Delete Old File
		if user.ImageProfile != "" {
			getNameOldFile := strings.Split(user.ImageProfile, "/")
			errDelete := os.RemoveAll(fmt.Sprintf("public/users/%v/profile_image/%v", helpers.GetLoginUserID(c), getNameOldFile[3]))
			if errDelete != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
			}
		}

		// Func Upload Image
		err = helpers.UploadImage(filePath, dirPath, *file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
		}

		user.ImageProfile = fileSrc
	}

	user.Name = input.Name

	if err := db.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Data: nil, Message: "Berhasil Update Profile User"})
}
