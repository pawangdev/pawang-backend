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

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type inputCategory struct {
	Name string `json:"name" form:"name" validate:"required"`
	Type string `json:"type" form:"type" validate:"required"`
}

func CategoryIndex(c echo.Context) error {
	db := config.ConnectDatabase()
	var categories []models.Category

	if err := db.Find(&categories).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Data: categories})
}

func CategoryShow(c echo.Context) error {
	db := config.ConnectDatabase()
	var category models.Category

	id := c.Param("categoryId")

	result := db.Preload("Category").Find(&category, "id = ?", id)

	if result.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, models.Response{Success: false, Message: "category not found", Data: nil})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Data: category})
}

func CategoryStore(c echo.Context) error {
	db := config.ConnectDatabase()
	var input inputCategory

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: err.Error(), Data: nil})
	}

	// Get Image from Input
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
	}

	// Setting Directory
	filePath := fmt.Sprintf("public/users/%v/categories/%v-%v", helpers.GetLoginUserID(c), time.Now().Unix(), file.Filename)
	fileSrc := fmt.Sprintf("categories/%v/categories/%v-%v", helpers.GetLoginUserID(c), time.Now().Unix(), file.Filename)
	dirPath := fmt.Sprintf("public/users/%v/categories/", helpers.GetLoginUserID(c))

	// Func Upload Image
	errUpload := helpers.UploadImage(filePath, dirPath, *file)
	if errUpload != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
	}

	category := new(models.Category)
	category.Name = input.Name
	category.Type = input.Type
	category.IconUrl = fileSrc
	category.UserID = helpers.GetLoginUserID(c)

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
	}

	return c.JSON(http.StatusCreated, models.Response{Success: true, Data: category, Message: "success create category"})
}

func CategoryUpdate(c echo.Context) error {
	db := config.ConnectDatabase()
	var input inputCategory
	category := new(models.Category)

	id := c.Param("categoryId")

	result := db.Find(&category, "id = ? AND user_id = ?", id, helpers.GetLoginUserID(c))

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Success: false, Message: "category not found", Data: nil})
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil})
	}

	if len(c.Request().MultipartForm.File) != 0 {
		file, err := c.FormFile("image")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
		}

		// Setting Directory
		filePath := fmt.Sprintf("public/users/%v/categories/%v-%v", helpers.GetLoginUserID(c), time.Now().Unix(), file.Filename)
		fileSrc := fmt.Sprintf("categories/%v/categories/%v-%v", helpers.GetLoginUserID(c), time.Now().Unix(), file.Filename)
		dirPath := fmt.Sprintf("public/users/%v/categories/", helpers.GetLoginUserID(c))

		// Delete Old File
		getNameOldFile := strings.Split(category.IconUrl, "/")
		errDelete := os.RemoveAll(fmt.Sprintf("public/users/%v/categories/%v", helpers.GetLoginUserID(c), getNameOldFile[3]))

		if errDelete != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
		}

		// Func Upload Image
		err = helpers.UploadImage(filePath, dirPath, *file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
		}

		category.IconUrl = fileSrc
	}

	category.Name = input.Name
	category.Type = input.Type
	category.UserID = helpers.GetLoginUserID(c)

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Message: "success update category", Data: category})
}

func CategoryDestroy(c echo.Context) error {
	db := config.ConnectDatabase()
	var category models.Category

	id := c.Param("categoryId")

	result := db.Find(&category, "id = ? AND user_id = ?", id, helpers.GetLoginUserID(c))

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "category not found", Data: nil})
	}

	// Delete Old File
	getNameOldFile := strings.Split(category.IconUrl, "/")
	errDelete := os.RemoveAll(fmt.Sprintf("public/users/%v/categories/%v", helpers.GetLoginUserID(c), getNameOldFile[3]))
	if errDelete != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: errDelete.Error(), Data: nil})
	}

	if err := db.Delete(&category, "id = ? AND user_id = ?", id, helpers.GetLoginUserID(c)).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Message: "success delete category", Data: nil})
}
