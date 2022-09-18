package response

import (
	"pawang-backend/entity"
	"time"
)

type CreateCategoryResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FormatCreateCategoryResponse(category entity.Category) CreateCategoryResponse {
	response := CreateCategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Icon:      category.Icon,
		Type:      category.Type,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	return response
}

type GetCategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Type string `json:"type"`
	// SubCategories []GetSubCategoryResponse `json:"sub_categories"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FormatGetCategoryResponse(category entity.Category) GetCategoryResponse {
	response := GetCategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Icon:      category.Icon,
		Type:      category.Type,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	// formatSubCategories := FormatGetSubCategoriesResponse(category.SubCategories)

	// response.SubCategories = formatSubCategories

	return response
}

func FormatGetCategoriesResponse(categories []entity.Category) []GetCategoryResponse {
	var categoriesResponse []GetCategoryResponse

	for _, category := range categories {
		categoryFormat := FormatGetCategoryResponse(category)
		categoriesResponse = append(categoriesResponse, categoryFormat)
	}

	return categoriesResponse
}
