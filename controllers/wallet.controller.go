package controllers

import (
	"net/http"
	"pawang-backend/config"
	"pawang-backend/helpers"
	"pawang-backend/models"

	"github.com/labstack/echo/v4"
)

type inputWallet struct {
	Name    string `json:"name" form:"name" validate:"required"`
	Balance uint64 `json:"balance" form:"balance" validate:"required,gte=0"`
}

func WalletIndex(c echo.Context) error {
	db := config.ConnectDatabase()
	var wallets []models.Wallet

	if err := db.Find(&wallets, "user_id = ?", helpers.GetLoginUserID(c)).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Data: wallets})
}

func WalletShow(c echo.Context) error {
	db := config.ConnectDatabase()
	var wallet models.Wallet

	id := c.Param("walletId")

	result := db.Where("id = ? AND user_id = ?", id, helpers.GetLoginUserID(c)).First(&wallet)

	if result.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, models.Response{Success: false, Message: "wallet not found", Data: nil})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Data: wallet})
}

func WalletStore(c echo.Context) error {
	db := config.ConnectDatabase()
	var input inputWallet

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	wallet := new(models.Wallet)
	wallet.Name = input.Name
	wallet.Balance = input.Balance
	wallet.UserID = helpers.GetLoginUserID(c)

	if err := db.Create(&wallet).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, models.Response{Success: true, Data: wallet, Message: "success create wallet"})
}

func WalletUpdate(c echo.Context) error {
	db := config.ConnectDatabase()
	var input inputWallet
	var wallet models.Wallet

	id := c.Param("walletId")

	result := db.Where("id = ? AND user_id = ?", id, helpers.GetLoginUserID(c)).First(&wallet)

	if result.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, models.Response{Success: false, Message: "wallet not found", Data: nil})
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: err.Error(), Data: nil})
	}

	wallet.Name = input.Name
	wallet.Balance = input.Balance
	wallet.UserID = helpers.GetLoginUserID(c)

	if err := db.Save(&wallet).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Message: "success update wallet", Data: wallet})
}

func WalletDestroy(c echo.Context) error {
	db := config.ConnectDatabase()
	var wallet models.Wallet

	id := c.Param("walletId")

	result := db.Where("id = ? AND user_id = ?", id, helpers.GetLoginUserID(c)).First(&wallet)
	if result.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, models.Response{Success: false, Message: "wallet not found", Data: nil})
	}

	if err := db.Delete(&wallet).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Success: false, Message: err.Error(), Data: nil})
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Message: "success delete wallet", Data: nil})
}
