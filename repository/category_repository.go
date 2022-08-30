package repository

import (
	"pawang-backend/entity"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(userID int, query string) ([]entity.Category, error)
	FindByID(userID int, categoryID int) (entity.Category, error)
	Insert(category entity.Category) (entity.Category, error)
	Update(category entity.Category) (entity.Category, error)
	Delete(category entity.Category) error
}

type categoryRepository struct {
	database *gorm.DB
}

func NewCategoryRepository(database *gorm.DB) *categoryRepository {
	return &categoryRepository{database}
}

func (repository *categoryRepository) FindAll(userID int, query string) ([]entity.Category, error) {
	var categories []entity.Category

	if len(query) == 0 {
		if err := repository.database.Preload("SubCategories", "user_id = ?", userID).Find(&categories).Error; err != nil {
			return categories, err
		}
	} else {
		if err := repository.database.Where("type = ?", query).Preload("SubCategories", "user_id = ?", userID).Find(&categories).Error; err != nil {
			return categories, err
		}
	}

	return categories, nil
}

func (repository *categoryRepository) FindByID(userID int, categoryID int) (entity.Category, error) {
	var category entity.Category

	if err := repository.database.Where("id = ?", categoryID).Preload("SubCategories", "user_id = ?", userID).Find(&category).Error; err != nil {
		return category, err
	}

	return category, nil
}

func (repository *categoryRepository) Insert(category entity.Category) (entity.Category, error) {
	if err := repository.database.Create(&category).Error; err != nil {
		return category, err
	}

	return category, nil
}

func (repository *categoryRepository) Update(category entity.Category) (entity.Category, error) {
	if err := repository.database.Save(&category).Error; err != nil {
		return category, err
	}

	return category, nil
}

func (repository *categoryRepository) Delete(category entity.Category) error {
	if err := repository.database.Delete(&category).Error; err != nil {
		return err
	}

	return nil
}
