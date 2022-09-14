package service

import (
	"errors"
	"pawang-backend/entity"
	"pawang-backend/helper"
	"pawang-backend/model/request"
	"pawang-backend/repository"
	"time"
)

type UserService interface {
	Register(input request.RegisterUserRequest) (entity.User, error)
	Login(input request.LoginUserInput) (entity.User, error)
	ChangePassword(userID int, input request.UserChangePasswordRequest) (entity.User, error)
	ChangeProfile(userID int, input request.UserChangeProfileRequest) (entity.User, error)
	Profile(userID int) (entity.User, error)
	RequestResetPasswordToken(input request.UserResetPasswordRequest) (entity.UserResetPassword, error)
	VerifyResetPasswordToken(input request.UserResetPasswordTokenRequest) (bool, error)
	ResetPasswordConfirmation(input request.UserResetPasswordConfirmation) (bool, error)
}

type userService struct {
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
}

func NewUserService(userRepository repository.UserRepository, walletRepository repository.WalletRepository) *userService {
	return &userService{userRepository, walletRepository}
}

func (service *userService) Register(input request.RegisterUserRequest) (entity.User, error) {
	user := entity.User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Phone = input.Phone
	user.Gender = input.Gender

	passwordHash, err := helper.HashPassword(input.Password)
	if err != nil {
		return user, err
	}

	user.Password = passwordHash

	newUser, err := service.userRepository.Insert(user)
	if err != nil {
		return user, err
	}

	wallet := entity.Wallet{}
	wallet.Name = "Dompet"
	wallet.Balance = 0
	wallet.UserID = newUser.ID
	_, err = service.walletRepository.Insert(wallet)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (service *userService) Login(input request.LoginUserInput) (entity.User, error) {
	email := input.Email
	password := input.Password

	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("the email does not exist")
	}

	err = helper.ComparePassword(user.Password, password)
	if err != nil {
		return user, errors.New("passwords mismatch")
	}

	return user, nil
}

func (service *userService) ChangePassword(userID int, input request.UserChangePasswordRequest) (entity.User, error) {
	oldPassword := input.OldPassword
	newPassword := input.NewPassword

	user, err := service.userRepository.FindByID(userID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("the email does not exist")
	}

	err = helper.ComparePassword(user.Password, oldPassword)
	if err != nil {
		return user, errors.New("passwords mismatch")
	}

	newPasswordHash, err := helper.HashPassword(newPassword)
	if err != nil {
		return user, err
	}

	user.Password = newPasswordHash

	updateUser, err := service.userRepository.Update(user)
	if err != nil {
		return updateUser, err
	}

	return updateUser, nil
}

func (service *userService) ChangeProfile(userID int, input request.UserChangeProfileRequest) (entity.User, error) {
	user, err := service.userRepository.FindByID(userID)
	if err != nil {
		return user, err
	}

	user.Name = input.Name
	user.Phone = input.Phone
	user.Gender = input.Gender

	updateUser, err := service.userRepository.Update(user)
	if err != nil {
		return updateUser, err
	}

	return updateUser, nil
}

func (service *userService) Profile(userID int) (entity.User, error) {
	user, err := service.userRepository.FindByID(userID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (service *userService) RequestResetPasswordToken(input request.UserResetPasswordRequest) (entity.UserResetPassword, error) {
	token := entity.UserResetPassword{}

	email := input.Email
	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return token, err
	}

	if user.ID == 0 {
		return token, errors.New("the email does not exist")
	}

	token.Email = user.Email
	token.Token = helper.GenerateRandomKey(8)
	token.ExpiredAt = time.Now().Add(time.Minute * 10)

	newToken, err := service.userRepository.InsertTokenResetPassword(token)
	if err != nil {
		return newToken, err
	}

	return newToken, nil
}

func (service *userService) VerifyResetPasswordToken(input request.UserResetPasswordTokenRequest) (bool, error) {
	token, err := service.userRepository.FindTokenResetPassword(input.Token)
	if err != nil {
		return false, errors.New("token does not same")
	}

	// user, err := service.userRepository.FindByEmail(token.Email)
	// if err != nil {
	// 	return user, err
	// }

	// if token.Token != input.Token {
	// 	return user, errors.New("token does not same")
	// }

	// if token.Email != user.Email {
	// 	return user, errors.New("token does not same")
	// }

	now := time.Now()

	if now.Before(token.ExpiredAt) {
		return true, nil
	} else {
		return false, errors.New("token has expired")
	}

	// updateUser, err := service.userRepository.Update(user)
	// if err != nil {
	// 	return updateUser, err
	// }
}

func (service *userService) ResetPasswordConfirmation(input request.UserResetPasswordConfirmation) (bool, error) {
	token, err := service.userRepository.FindTokenResetPassword(input.Token)
	if err != nil {
		return false, errors.New("token does not same")
	}

	user, err := service.userRepository.FindByEmail(token.Email)
	if err != nil {
		return false, err
	}

	if token.Token != input.Token {
		return false, errors.New("token does not same")
	}

	if token.Email != user.Email {
		return false, errors.New("token does not same")
	}

	newPasswordHash, err := helper.HashPassword(input.Password)
	if err != nil {
		return false, err
	}
	user.Password = newPasswordHash

	_, err = service.userRepository.Update(user)
	if err != nil {
		return false, err
	}

	err = service.userRepository.DeleteTokenResetPassword(token)
	if err != nil {
		return false, err
	}

	return true, nil
}
