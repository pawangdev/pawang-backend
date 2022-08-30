package handler

import (
	"pawang-backend/exception"
	"pawang-backend/model/request"
	"pawang-backend/model/response"
	"pawang-backend/service"

	"github.com/gofiber/fiber/v2"
)

type subCategoryHandler struct {
	subCategoryService service.SubCategoryService
	authService        service.AuthService
}

func NewSubCategoryHandler(subCategoryService service.SubCategoryService, authService service.AuthService) *subCategoryHandler {
	return &subCategoryHandler{subCategoryService, authService}
}

func (handler *subCategoryHandler) CreateSubCategory(c *fiber.Ctx) error {
	var input request.CreateSubCategoryRequest

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

	newSubCategory, err := handler.subCategoryService.CreateSubCategory(userID, input)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formatter := response.FormatGetSubCategoryResponse(newSubCategory)
	response := response.ResponseSuccess(true, "Subcategory created successfully", formatter)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (handler *subCategoryHandler) UpdateSubCategory(c *fiber.Ctx) error {
	// Get SubCategory ID
	subCategoryID, err := c.ParamsInt("subcategoryId")
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var input request.UpdateSubCategory

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

	updateSubCategory, err := handler.subCategoryService.UpdateSubCategory(subCategoryID, userID, input)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formatter := response.FormatGetSubCategoryResponse(updateSubCategory)
	response := response.ResponseSuccess(true, "Subcategory updated successfully", formatter)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (handler *subCategoryHandler) DeleteSubCategory(c *fiber.Ctx) error {
	// Get SubCategory ID
	subCategoryID, err := c.ParamsInt("subcategoryId")
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Get User ID Current User Logged
	userID := handler.authService.CurrentLoggedUserID(c)

	err = handler.subCategoryService.DeleteSubCategory(subCategoryID, userID)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := response.ResponseSuccess(true, "Subcategory deleted successfully", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
