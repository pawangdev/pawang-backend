package request

import (
	"time"
)

type CreateTransactionRequest struct {
	Amount        int    `json:"amount" form:"amount" xml:"amount" validate:"required"`
	CategoryID    int    `json:"category_id" form:"category_id" xml:"category_id" validate:"required"`
	SubCategoryID int    `json:"sub_category_id" form:"sub_category_id" xml:"sub_category_id"`
	WalletID      int    `json:"wallet_id" form:"wallet_id" xml:"wallet_id" validate:"required"`
	Type          string `json:"type" form:"type" xml:"type" validate:"required,oneof=income outcome"`
	Description   string `json:"description" form:"description" xml:"description"`
	ImageUrl      string
	Date          time.Time `json:"date" form:"date" xml:"date" validate:"required"`
}
