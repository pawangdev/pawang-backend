package handler

import (
	"fmt"
	"pawang-backend/exception"
	"pawang-backend/model/request"
	"pawang-backend/model/response"
	"pawang-backend/service"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userService  service.UserService
	authService  service.AuthService
	emailService service.EmailService
}

func NewUserHandler(userService service.UserService, authService service.AuthService, emailService service.EmailService) *userHandler {
	return &userHandler{userService, authService, emailService}
}

func (handler *userHandler) RegisterUser(c *fiber.Ctx) error {
	var input request.RegisterUserRequest

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

	newUser, err := handler.userService.Register(input)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	token, err := handler.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formatter := response.FormatRegisterUserResponse(newUser, token)
	response := response.ResponseSuccess(true, "Register new user successfully", formatter)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (handler *userHandler) LoginUser(c *fiber.Ctx) error {
	var input request.LoginUserInput

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

	user, err := handler.userService.Login(input)
	if err != nil {
		response := exception.ResponseError(false, err.Error(), fiber.StatusUnauthorized, nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	token, err := handler.authService.GenerateToken(user.ID)
	if err != nil {
		response := exception.ResponseError(false, err.Error(), fiber.StatusUnauthorized, nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	formatter := response.FormatRegisterUserResponse(user, token)
	response := response.ResponseSuccess(true, "User logged in successfully", formatter)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (handler *userHandler) ChangePassword(c *fiber.Ctx) error {
	var input request.UserChangePasswordRequest

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

	_, err := handler.userService.ChangePassword(userID, input)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := response.ResponseSuccess(true, "Change password has successfully", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (handler *userHandler) ChangeProfile(c *fiber.Ctx) error {
	var input request.UserChangeProfileRequest

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

	user, err := handler.userService.ChangeProfile(userID, input)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formatter := response.FormatUserProfileResponse(user)
	response := response.ResponseSuccess(true, "Change profile has successfully", formatter)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (handler *userHandler) UserProfile(c *fiber.Ctx) error {
	// Get User ID Current User Logged
	userID := handler.authService.CurrentLoggedUserID(c)

	user, err := handler.userService.Profile(userID)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formatter := response.FormatUserProfileResponse(user)
	response := response.ResponseSuccess(true, "Getting profile successfully", formatter)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (handler *userHandler) RequestResetPasswordToken(c *fiber.Ctx) error {
	var input request.UserResetPasswordRequest

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

	tokenReset, err := handler.userService.RequestResetPasswordToken(input)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	bodyMessage := fmt.Sprintf("<p>Gunakan kode ini untuk mengatur ulang kata sandi akun Anda: <strong>%s</strong>. Kode hanya berlaku 10 menit.</p>", tokenReset.Token)
	err = handler.emailService.SendEmail(input.Email, "Kode Lupa Kata Sandi", bodyMessage)
	if err != nil {
		fmt.Println(err.Error())
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := response.ResponseSuccess(true, "Email has been sent!", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (handler *userHandler) VerifyResetPasswordToken(c *fiber.Ctx) error {
	var input request.UserResetPasswordTokenRequest

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

	_, err := handler.userService.VerifyResetPasswordToken(input)
	if err != nil {
		response := exception.ResponseError(false, "", fiber.StatusBadRequest, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := response.ResponseSuccess(true, "Success reset password, please re-login !", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
