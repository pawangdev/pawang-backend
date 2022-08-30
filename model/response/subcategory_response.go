package response

import (
	"pawang-backend/entity"
	"time"
)

type GetSubCategoryResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	CategoryID int       `json:"category_id"`
	UserID     int       `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func FormatGetSubCategoryResponse(subCategory entity.SubCategory) GetSubCategoryResponse {
	response := GetSubCategoryResponse{}
	response.ID = subCategory.ID
	response.Name = subCategory.Name
	response.CategoryID = subCategory.CategoryID
	response.UserID = subCategory.UserID
	response.CreatedAt = subCategory.CreatedAt
	response.UpdatedAt = subCategory.UpdatedAt

	return response
}

func FormatGetSubCategoriesResponse(subCategory []entity.SubCategory) []GetSubCategoryResponse {
	var data []GetSubCategoryResponse

	for _, subcategory := range subCategory {
		format := FormatGetSubCategoryResponse(subcategory)
		data = append(data, format)
	}

	return data
}
