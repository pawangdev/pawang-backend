package repository

import (
	"pawang-backend/entity"

	"gorm.io/gorm"
)

type SubCategoryRepository interface {
	FindByID(subcategoryID int) (entity.SubCategory, error)
	Insert(subcategory entity.SubCategory) (entity.SubCategory, error)
	Update(subcategory entity.SubCategory) (entity.SubCategory, error)
	Delete(subcategory entity.SubCategory) error
}

type subCategoryRepository struct {
	database *gorm.DB
}

func NewSubCategoryRepository(database *gorm.DB) *subCategoryRepository {
	return &subCategoryRepository{database}
}

func (repository *subCategoryRepository) FindByID(subcategoryID int) (entity.SubCategory, error) {
	var subcategory entity.SubCategory

	if err := repository.database.Where("id = ?", subcategoryID).Find(&subcategory).Error; err != nil {
		return subcategory, err
	}

	return subcategory, nil
}

func (repository *subCategoryRepository) Insert(subcategory entity.SubCategory) (entity.SubCategory, error) {
	if err := repository.database.Create(&subcategory).Error; err != nil {
		return subcategory, err
	}

	return subcategory, nil
}

func (repository *subCategoryRepository) Update(subcategory entity.SubCategory) (entity.SubCategory, error) {
	if err := repository.database.Save(&subcategory).Error; err != nil {
		return subcategory, err
	}

	return subcategory, nil
}

func (repository *subCategoryRepository) Delete(subcategory entity.SubCategory) error {
	if err := repository.database.Delete(&subcategory).Error; err != nil {
		return err
	}

	return nil
}
