package handler

import (
	"fmt"
	"pawang-backend/exception"
	"pawang-backend/helper"
	"pawang-backend/model/request"
	"pawang-backend/model/response"
	"pawang-backend/service"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type transactionHandler struct {
	transactionService service.TransactionService
	authService        service.AuthService
}

func NewTransactionHandler(transactionService service.TransactionService, authService service.AuthService) *transactionHandler {
	return &transactionHandler{transactionService, authService}
}

func (handler *transactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var input request.CreateTransactionRequest

	if err := c.BodyParser(&input); err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validate Input
	if errors := exception.ValidateInput(input); errors != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Get User ID Current User Logged
	userID := handler.authService.CurrentLoggedUserID(c)

	// file, _ := c.FormFile("image")

	isUpload := c.FormValue("is_upload", "false")

	if isUpload == "true" {
		file, err := c.FormFile("image")
		if err != nil {
			response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		dirPath := fmt.Sprintf("./public/uploads/transactions/%v", userID)
		fileName := fmt.Sprintf("./public/uploads/transactions/%v/%d-%s", userID, time.Now().Unix(), file.Filename)
		imageUrl := fmt.Sprintf("/api/storage/uploads/transactions/%v/%d-%s", userID, time.Now().Unix(), file.Filename)

		err = helper.CreateFolder(dirPath)
		if err != nil {
			response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		err = c.SaveFile(file, fileName)
		if err != nil {
			response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		input.ImageUrl = imageUrl
	}

	newTransaction, err := handler.transactionService.CreateTransaction(userID, input)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formatter := response.FormatCreateTransactionResponse(newTransaction)
	response := response.ResponseSuccess(true, "Transaction created successfully", formatter)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (handler *transactionHandler) GetTransactions(c *fiber.Ctx) error {
	// Get User ID Current User Logged
	userID := handler.authService.CurrentLoggedUserID(c)

	transactions, err := handler.transactionService.GetTransactions(userID)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formatter := response.FormatGetTransactionsResponse(transactions)
	response := response.ResponseSuccess(true, "Successfully get transaction", formatter)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (handler *transactionHandler) UpdateTransaction(c *fiber.Ctx) error {
	transactionID, err := c.ParamsInt("transactionId")
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var input request.CreateTransactionRequest
	if err := c.BodyParser(&input); err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validation Input
	if errors := exception.ValidateInput(input); errors != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Get User ID Current User Logged
	userID := handler.authService.CurrentLoggedUserID(c)

	// Get Old Transaction
	transaction, err := handler.transactionService.GetTransactionByID(transactionID, userID)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusNotFound, err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	isUpload := c.FormValue("is_upload", "false")

	if isUpload == "true" {
		file, err := c.FormFile("image")
		if err != nil {
			response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		dirPath := fmt.Sprintf("./public/uploads/transactions/%v", userID)
		fileName := fmt.Sprintf("./public/uploads/transactions/%v/%d-%s", userID, time.Now().Unix(), file.Filename)
		imageUrl := fmt.Sprintf("/api/storage/uploads/transactions/%v/%d-%s", userID, time.Now().Unix(), file.Filename)

		err = helper.CreateFolder(dirPath)
		if err != nil {
			response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		err = c.SaveFile(file, fileName)
		if err != nil {
			response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		// Delete Old File
		if transaction.ImageUrl != "" {
			getNameOldFile := strings.Split(transaction.ImageUrl, "/")[7]
			fileOldPath := fmt.Sprintf("./public/uploads/transactions/%v/%v", userID, getNameOldFile)
			err = helper.DeleteFile(fileOldPath)
			if err != nil {
				response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}
		}

		input.ImageUrl = imageUrl
	}

	updateTransaction, err := handler.transactionService.UpdateTransaction(transactionID, userID, input)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusNotFound, err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	formatter := response.FormatCreateTransactionResponse(updateTransaction)
	response := response.ResponseSuccess(true, "Transaction updated successfully", formatter)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (handler *transactionHandler) DeleteTransaction(c *fiber.Ctx) error {
	transactionID, err := c.ParamsInt("transactionId")
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Get User ID Current User Logged
	userID := handler.authService.CurrentLoggedUserID(c)

	// Get Old Transaction
	transaction, err := handler.transactionService.GetTransactionByID(transactionID, userID)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusNotFound, err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Delete Old File
	if transaction.ImageUrl != "" {
		getNameOldFile := strings.Split(transaction.ImageUrl, "/")[7]
		fileOldPath := fmt.Sprintf("./public/uploads/transactions/%v/%v", userID, getNameOldFile)
		err = helper.DeleteFile(fileOldPath)
		if err != nil {
			response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	err = handler.transactionService.DeleteTransaction(transactionID, userID)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := response.ResponseSuccess(true, "Transaction deleted successfully", nil)
	return c.Status(fiber.StatusOK).JSON(response)

}
