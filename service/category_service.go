package service

import (
	"pawang-backend/entity"
	"pawang-backend/repository"
)

type CategoryService interface {
	GetCategories(userID int, query string) ([]entity.Category, error)
	// CreateCategory(userID int, input request.CreateCategoryRequest) (entity.Category, error)
	// UpdateCategory(categoryID int, userID int, input request.CreateCategoryRequest) (entity.Category, error)
	// DeleteCategory(categoryID int, userID int) error
}

type categoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(repository repository.CategoryRepository) *categoryService {
	return &categoryService{repository}
}

func (service *categoryService) GetCategories(userID int, query string) ([]entity.Category, error) {
	categories, err := service.repository.FindAll(userID, query)
	if err != nil {
		return categories, err
	}

	return categories, nil
}

// func (service *categoryService) CreateCategory(userID int, input request.CreateCategoryRequest) (entity.Category, error) {
// 	category := entity.Category{}
// 	category.Name = input.Name
// 	category.Type = input.Type
// 	category.Icon = input.Icon

// 	newCategory, err := service.repository.InsertSubCategory(category)
// 	if err != nil {
// 		return newCategory, err
// 	}

// 	return newCategory, nil
// }

// func (service *categoryService) UpdateCategory(categoryID int, userID int, input request.CreateCategoryRequest) (entity.Category, error) {
// 	category, err := service.repository.FindByID(categoryID)
// 	if err != nil {
// 		return category, err
// 	}

// 	if category.ID == 0 {
// 		return category, errors.New("the category does not exist")
// 	}

// 	category.Name = input.Name
// 	category.Type = input.Type
// 	category.Icon = input.Icon

// 	updateCategory, err := service.repository.UpdateSubCategory(category)
// 	if err != nil {
// 		return updateCategory, err
// 	}

// 	return updateCategory, nil
// }

// func (service *categoryService) DeleteCategory(categoryID int, userID int) error {
// 	category, err := service.repository.FindByID(categoryID)
// 	if err != nil {
// 		return err
// 	}

// 	if category.ID == 0 {
// 		return errors.New("the category does not exist")
// 	}

// 	err = service.repository.DeleteSubCategory(category)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
