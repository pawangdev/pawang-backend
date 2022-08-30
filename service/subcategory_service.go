package service

import (
	"errors"
	"pawang-backend/entity"
	"pawang-backend/model/request"
	"pawang-backend/repository"
)

type SubCategoryService interface {
	CreateSubCategory(userID int, input request.CreateSubCategoryRequest) (entity.SubCategory, error)
	UpdateSubCategory(subCategoryID int, userID int, input request.UpdateSubCategory) (entity.SubCategory, error)
	DeleteSubCategory(subCategoryID int, userID int) error
}

type subCategoryService struct {
	repository repository.SubCategoryRepository
}

func NewSubCategoryService(repository repository.SubCategoryRepository) *subCategoryService {
	return &subCategoryService{repository}
}

func (service *subCategoryService) CreateSubCategory(userID int, input request.CreateSubCategoryRequest) (entity.SubCategory, error) {
	subCategory := entity.SubCategory{}
	subCategory.Name = input.Name
	subCategory.CategoryID = input.CategoryID
	subCategory.UserID = userID

	newSubCategory, err := service.repository.Insert(subCategory)
	if err != nil {
		return newSubCategory, err
	}

	return newSubCategory, nil
}

func (service *subCategoryService) UpdateSubCategory(subCategoryID int, userID int, input request.UpdateSubCategory) (entity.SubCategory, error) {
	subCategory, err := service.repository.FindByID(subCategoryID)
	if err != nil {
		return subCategory, err
	}

	if subCategory.ID == 0 {
		return subCategory, errors.New("the subcategory does not exist")
	}

	if subCategory.UserID != userID {
		return subCategory, errors.New("the subcategory does not exist")
	}

	subCategory.Name = input.Name

	updateSubCategory, err := service.repository.Update(subCategory)
	if err != nil {
		return updateSubCategory, err
	}

	return updateSubCategory, nil
}

func (service *subCategoryService) DeleteSubCategory(subCategoryID int, userID int) error {
	subCategory, err := service.repository.FindByID(subCategoryID)
	if err != nil {
		return err
	}

	if subCategory.ID == 0 {
		return errors.New("the subcategory does not exist")
	}

	if subCategory.UserID != userID {
		return errors.New("the subcategory does not exist")
	}

	err = service.repository.Delete(subCategory)
	if err != nil {
		return err
	}

	return nil
}
