package handler

import (
	"pawang-backend/exception"
	"pawang-backend/model/request"
	"pawang-backend/model/response"
	"pawang-backend/service"

	"github.com/gofiber/fiber/v2"
)

type walletHandler struct {
	walletService service.WalletService
	authService   service.AuthService
}

func NewWalletHandler(walletService service.WalletService, authService service.AuthService) *walletHandler {
	return &walletHandler{walletService, authService}
}

func (handler *walletHandler) GetWallets(c *fiber.Ctx) error {
	// Get User ID Current User Logged
	userID := handler.authService.CurrentLoggedUserID(c)

	wallets, err := handler.walletService.GetWallets(userID)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formatter := response.FormatGetWalletsResponse(wallets)
	response := response.ResponseSuccess(true, "Successfully get wallets", formatter)
	return c.Status(fiber.StatusOK).JSON(response)

}

func (handler *walletHandler) CreateWallet(c *fiber.Ctx) error {
	var input request.CreateWalletRequest

	if err := c.BodyParser(&input); err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validation Input
	errors := exception.ValidateInput(input)
	if errors != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Get User ID Current User Logged
	userID := handler.authService.CurrentLoggedUserID(c)

	newWallet, err := handler.walletService.CreateWallet(userID, input)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formatter := response.FormatCreateWalletResponse(newWallet)
	response := response.ResponseSuccess(true, "Wallet created successfully", formatter)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (handler *walletHandler) UpdateWallet(c *fiber.Ctx) error {
	walletID, err := c.ParamsInt("walletId")
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var input request.UpdateWalletRequest
	if err := c.BodyParser(&input); err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validation Input
	errors := exception.ValidateInput(input)
	if errors != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Get User ID Current User Logged
	userID := handler.authService.CurrentLoggedUserID(c)

	updateWallet, err := handler.walletService.UpdateWallet(walletID, userID, input)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formatter := response.FormatCreateWalletResponse(updateWallet)
	response := response.ResponseSuccess(true, "Wallet updated successfully", formatter)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (handler *walletHandler) DeleteWallet(c *fiber.Ctx) error {
	walletID, err := c.ParamsInt("walletId")
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Get User ID Current User Logged
	userID := handler.authService.CurrentLoggedUserID(c)

	err = handler.walletService.DeleteWallet(walletID, userID)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := response.ResponseSuccess(true, "Wallet deleted successfully", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
