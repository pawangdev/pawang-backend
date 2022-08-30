package request

type CreateWalletRequest struct {
	Name    string `json:"name" form:"name" xml:"name" validate:"required"`
	Balance int    `json:"balance" form:"balance" xml:"balance" validate:"number"`
}

type UpdateWalletRequest struct {
	Name    string `json:"name" form:"name" xml:"name" validate:"required"`
	Balance int    `json:"balance" form:"balance" xml:"balance" validate:"number"`
}
