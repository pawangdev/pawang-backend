package repository

import (
	"pawang-backend/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserNotificationRepository interface {
	Create(notification entity.UserNotification) (entity.UserNotification, error)
	Update(notification entity.UserNotification) (entity.UserNotification, error)
	FindByOneSignalID(onesignalId string) (entity.UserNotification, error)
	Delete(notification entity.UserNotification) error
}

type userNotificationRepository struct {
	database *gorm.DB
}

func NewUserNotificationRepository(database *gorm.DB) *userNotificationRepository {
	return &userNotificationRepository{database}
}

func (repository *userNotificationRepository) Create(notification entity.UserNotification) (entity.UserNotification, error) {
	if err := repository.database.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&notification).Error; err != nil {
		return notification, err
	}

	return notification, nil
}

func (repository *userNotificationRepository) Update(notification entity.UserNotification) (entity.UserNotification, error) {
	if err := repository.database.Save(&notification).Error; err != nil {
		return notification, err
	}

	return notification, nil
}

func (repository *userNotificationRepository) FindByOneSignalID(onesignalId string) (entity.UserNotification, error) {
	notification := entity.UserNotification{}

	if err := repository.database.Where("onesignal_id = ?", onesignalId).Find(&notification).Error; err != nil {
		return notification, err
	}

	return notification, nil
}

func (repository *userNotificationRepository) Delete(notification entity.UserNotification) error {
	if err := repository.database.Delete(&notification).Error; err != nil {
		return err
	}

	return nil
}
