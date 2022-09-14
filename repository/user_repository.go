package repository

import (
	"pawang-backend/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Insert(user entity.User) (entity.User, error)
	FindByEmail(email string) (entity.User, error)
	FindByID(userID int) (entity.User, error)
	Update(user entity.User) (entity.User, error)
	InsertTokenResetPassword(token entity.UserResetPassword) (entity.UserResetPassword, error)
	FindTokenResetPassword(token string) (entity.UserResetPassword, error)
	DeleteTokenResetPassword(token entity.UserResetPassword) error
}

type userRepository struct {
	database *gorm.DB
}

func NewUserRepository(database *gorm.DB) *userRepository {
	return &userRepository{database}
}

func (repository *userRepository) Insert(user entity.User) (entity.User, error) {
	if err := repository.database.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (repository *userRepository) FindByEmail(email string) (entity.User, error) {
	var user entity.User

	if err := repository.database.Where("email = ?", email).Find(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (repository *userRepository) FindByID(userID int) (entity.User, error) {
	var user entity.User

	if err := repository.database.Where("id = ?", userID).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (repository *userRepository) Update(user entity.User) (entity.User, error) {
	if err := repository.database.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (repository *userRepository) InsertTokenResetPassword(token entity.UserResetPassword) (entity.UserResetPassword, error) {
	if err := repository.database.Create(&token).Error; err != nil {
		return token, err
	}

	return token, nil
}

func (repository *userRepository) FindTokenResetPassword(token string) (entity.UserResetPassword, error) {
	var tokens entity.UserResetPassword

	if err := repository.database.Where("token = ?", token).First(&tokens).Error; err != nil {
		return tokens, err
	}

	return tokens, nil
}

func (repository *userRepository) DeleteTokenResetPassword(token entity.UserResetPassword) error {
	if err := repository.database.Delete(&token).Error; err != nil {
		return err
	}

	return nil
}
