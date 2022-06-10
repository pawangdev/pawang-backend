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
)

type inputTransaction struct {
	Amount      uint64    `json:"amount" form:"amount" validate:"required,gte=0"`
	CategoryID  uint      `json:"category_id" form:"category_id" validate:"required"`
	WalletID    uint      `json:"wallet_id" form:"wallet_id" validate:"required"`
	Type        string    `json:"type" form:"type" validate:"required"`
	Description string    `json:"description" form:"description"`
	ImageUrl    string    `json:"image_url" form:"image_url"`
	Date        time.Time `json:"date" form:"date" validate:"required"`
}

func TransactionIndex(c echo.Context) error {
	db := config.ConnectDatabase()
	var transactions []models.Transaction

	if err := db.Preload("Wallet").Preload("Category").Find(&transactions, "user_id = ?", helpers.GetLoginUserID(c)).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Data: transactions})
}

func TransactionShow(c echo.Context) error {
	db := config.ConnectDatabase()
	transaction := new(models.Transaction)

	id := c.Param("transactionId")

	result := db.Preload("Wallet").Preload("Category").Where("id = ? AND user_id = ?", id, helpers.GetLoginUserID(c)).First(&transaction)

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Success: false, Data: nil, Message: "transaction not found"})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Data: transaction})
}

func TransactionStore(c echo.Context) error {
	db := config.ConnectDatabase()
	var input inputTransaction
	transaction := new(models.Transaction)
	wallet := new(models.Wallet)

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	getWallet := db.Where("id = ?", input.WalletID).First(&wallet)

	if getWallet.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Success: false, Data: nil, Message: "wallet not found"})
	}

	if input.Type == "outcome" {
		if wallet.Balance < input.Amount {
			return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil, Message: "the balance is not sufficient"})
		}

		wallet.Balance = wallet.Balance - input.Amount
	} else if input.Type == "income" {
		wallet.Balance = wallet.Balance + input.Amount
	} else {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil, Message: "type not found"})
	}

	if len(c.Request().MultipartForm.File) != 0 {
		// Get Form Image
		file, err := c.FormFile("image")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
		}

		// Setting Directory
		filePath := fmt.Sprintf("public/users/%v/transactions/%v-%v", helpers.GetLoginUserID(c), time.Now().Unix(), file.Filename)
		fileSrc := fmt.Sprintf("transactions/%v/transactions/%v-%v", helpers.GetLoginUserID(c), time.Now().Unix(), file.Filename)
		dirPath := fmt.Sprintf("public/users/%v/transactions/", helpers.GetLoginUserID(c))

		// Func Upload Image
		err = helpers.UploadImage(filePath, dirPath, *file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
		}

		transaction.ImageUrl = fileSrc
	}

	transaction.Amount = input.Amount
	transaction.CategoryID = input.CategoryID
	transaction.WalletID = input.WalletID
	transaction.Type = input.Type
	transaction.Description = input.Description
	transaction.Date = input.Date
	transaction.UserID = helpers.GetLoginUserID(c)

	if err := db.Create(&transaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	if err := db.Save(&wallet).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, models.Response{Success: true, Data: transaction, Message: "success create transaction"})
}

func TransactionUpdate(c echo.Context) error {
	db := config.ConnectDatabase()
	var input inputTransaction
	transaction := new(models.Transaction)
	wallet := new(models.Wallet)

	id := c.Param("transactionId")

	result := db.Where("id = ? AND user_id = ?", id, helpers.GetLoginUserID(c)).First(&transaction)

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Success: false, Data: nil, Message: "transaction not found"})
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	getWallet := db.Where("id = ?", input.WalletID).First(&wallet)

	if getWallet.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Success: false, Data: nil, Message: "wallet not found"})
	}

	if input.Type == "outcome" {
		if wallet.Balance < input.Amount {
			return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil, Message: "the balance is not sufficient"})
		}

		wallet.Balance = (wallet.Balance - transaction.Amount) - input.Amount
	} else if input.Type == "income" {
		wallet.Balance = (wallet.Balance - transaction.Amount) + input.Amount
	} else {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil, Message: "type not found"})
	}

	if len(c.Request().MultipartForm.File) != 0 {
		// Get Form Image
		file, err := c.FormFile("image")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
		}

		// Setting Directory
		filePath := fmt.Sprintf("public/users/%v/transactions/%v-%v", helpers.GetLoginUserID(c), time.Now().Unix(), file.Filename)
		fileSrc := fmt.Sprintf("transactions/%v/transactions/%v-%v", helpers.GetLoginUserID(c), time.Now().Unix(), file.Filename)
		dirPath := fmt.Sprintf("public/users/%v/transactions/", helpers.GetLoginUserID(c))

		// Delete Old File
		getNameOldFile := strings.Split(transaction.ImageUrl, "/")
		errDelete := os.RemoveAll(fmt.Sprintf("public/users/%v/transactions/%v", helpers.GetLoginUserID(c), getNameOldFile[3]))
		if errDelete != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
		}

		// Func Upload Image
		err = helpers.UploadImage(filePath, dirPath, *file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
		}

		transaction.ImageUrl = fileSrc
	}

	transaction.Amount = input.Amount
	transaction.CategoryID = input.CategoryID
	transaction.WalletID = input.WalletID
	transaction.Type = input.Type
	transaction.Description = input.Description
	transaction.Date = input.Date
	transaction.UserID = helpers.GetLoginUserID(c)

	if err := db.Save(&transaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	if err := db.Save(&wallet).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Data: transaction, Message: "success update transaction"})
}

func TransactionDestroy(c echo.Context) error {
	db := config.ConnectDatabase()
	transaction := new(models.Transaction)
	wallet := new(models.Wallet)

	id := c.Param("transactionId")

	result := db.Where("id = ? AND user_id = ?", id, helpers.GetLoginUserID(c)).First(&transaction)

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Success: false, Data: nil, Message: "transaction not found"})
	}

	// Restore Balance on Wallets
	getWallet := db.Where("id = ?", transaction.WalletID).First(&wallet)

	if getWallet.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Success: false, Data: nil, Message: "wallet not found"})
	}

	wallet.Balance = wallet.Balance + transaction.Amount

	if err := db.Save(&transaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	if err := db.Save(&wallet).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	// Delete Old File
	getNameOldFile := strings.Split(transaction.ImageUrl, "/")
	errDelete := os.RemoveAll(fmt.Sprintf("public/users/%v/transactions/%v", helpers.GetLoginUserID(c), getNameOldFile[3]))
	if errDelete != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: errDelete.Error(), Data: nil})
	}

	if err := db.Delete(&transaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Data: nil, Message: "success delete transaction"})
}
