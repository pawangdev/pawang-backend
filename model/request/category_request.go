package request

type CreateCategoryRequest struct {
	Name string `json:"name" form:"name" xml:"name" validate:"required"`
	Icon string `json:"icon" form:"icon" xml:"icon" validate:"required"`
	Type string `json:"type" form:"type" xml:"type" validate:"required,oneof=income outcome"`
}
