package request

type CreateSubCategoryRequest struct {
	Name       string `json:"name" form:"name" xml:"name" validate:"required"`
	CategoryID int    `json:"category_id" form:"category_id" xml:"category_id" validate:"required"`
}

type UpdateSubCategory struct {
	Name       string `json:"name" form:"name" xml:"name" validate:"required"`
}
