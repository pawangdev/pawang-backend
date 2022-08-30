package handler

import (
	"pawang-backend/exception"
	"pawang-backend/model/response"
	"pawang-backend/service"

	"github.com/gofiber/fiber/v2"
)

type categoryHandler struct {
	categoryService service.CategoryService
	authService     service.AuthService
}

func NewCategoryHandler(categoryService service.CategoryService, authService service.AuthService) *categoryHandler {
	return &categoryHandler{categoryService, authService}
}

// func (handler *categoryHandler) CreateCategory(c *fiber.Ctx) error {
// 	var input request.CreateCategoryRequest

// 	if err := c.BodyParser(&input); err != nil {
// 		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	// Validation Input
// 	errors := exception.ValidateInput(input)
// 	if errors != nil {
// 		response := exception.ResponseError(false, "", fiber.StatusBadRequest, errors)
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	// Get User ID Current User Logged
// 	userID := handler.authService.CurrentLoggedUserID(c)

// 	newCategory, err := handler.categoryService.CreateCategory(userID, input)
// 	if err != nil {
// 		response := exception.ResponseError(false, "", fiber.StatusBadRequest, errors)
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	formatter := response.FormatGetCategoryResponse(newCategory)
// 	response := response.ResponseSuccess(true, "Category created successfully", formatter)
// 	return c.Status(fiber.StatusCreated).JSON(response)
// }

func (handler *categoryHandler) GetCategories(c *fiber.Ctx) error {
	queryType := c.Query("type")
	// Get User ID Current User Logged
	userID := handler.authService.CurrentLoggedUserID(c)

	categories, err := handler.categoryService.GetCategories(userID, queryType)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formatter := response.FormatGetCategoriesResponse(categories)
	response := response.ResponseSuccess(true, "Successfully get categories", formatter)
	return c.Status(fiber.StatusOK).JSON(response)
}

// func (handler *categoryHandler) UpdateCategory(c *fiber.Ctx) error {
// 	var input request.CreateCategoryRequest

// 	categoryID, err := c.ParamsInt("categoryId")
// 	if err != nil {
// 		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	if err := c.BodyParser(&input); err != nil {
// 		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	// Validation Input
// 	errors := exception.ValidateInput(input)
// 	if errors != nil {
// 		response := exception.ResponseError(false, "", fiber.StatusBadRequest, errors)
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	// Get User ID Current User Logged
// 	userID := handler.authService.CurrentLoggedUserID(c)

// 	updateCategory, err := handler.categoryService.UpdateCategory(categoryID, userID, input)
// 	if err != nil {
// 		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	formatter := response.FormatGetCategoryResponse(updateCategory)
// 	response := response.ResponseSuccess(true, "Category updated successfully", formatter)
// 	return c.Status(fiber.StatusOK).JSON(response)
// }

// func (handler *categoryHandler) DeleteCategory(c *fiber.Ctx) error {
// 	categoryID, err := c.ParamsInt("categoryId")
// 	if err != nil {
// 		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	// Get User ID Current User Logged
// 	userID := handler.authService.CurrentLoggedUserID(c)

// 	err = handler.categoryService.DeleteCategory(categoryID, userID)
// 	if err != nil {
// 		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(response)
// 	}

// 	response := response.ResponseSuccess(true, "Category deleted successfully", nil)
// 	return c.Status(fiber.StatusOK).JSON(response)

// }
