package response

import (
	"pawang-backend/entity"
	"time"
)

type CreateTransactionResponse struct {
	ID         int `json:"id"`
	Amount     int `json:"amount"`
	CategoryID int `json:"category_id"`
	// SubCategoryID int       `json:"subcategory_id"`
	WalletID    int       `json:"wallet_id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image_url"`
	Date        time.Time `json:"date"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func FormatCreateTransactionResponse(transaction entity.Transaction) CreateTransactionResponse {
	response := CreateTransactionResponse{
		ID:         transaction.ID,
		Amount:     transaction.Amount,
		CategoryID: transaction.CategoryID,
		// SubCategoryID: transaction.SubCategoryID,
		WalletID:    transaction.WalletID,
		Type:        transaction.Type,
		Description: transaction.Description,
		ImageUrl:    transaction.ImageUrl,
		Date:        transaction.Date,
		UserID:      transaction.UserID,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}

	return response
}

type GetTransactionResponse struct {
	ID         int `json:"id"`
	Amount     int `json:"amount"`
	CategoryID int `json:"category_id"`
	// SubCategoryID int         `json:"subcategory_id"`
	WalletID    int         `json:"wallet_id"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	ImageUrl    string      `json:"image_url"`
	Date        time.Time   `json:"date"`
	UserID      int         `json:"user_id"`
	Category    interface{} `json:"category"`
	SubCategory interface{} `json:"sub_category"`
	Wallet      interface{} `json:"wallet"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func FormatGetTransactionResponse(transaction entity.Transaction) GetTransactionResponse {
	response := GetTransactionResponse{
		ID:         transaction.ID,
		Amount:     transaction.Amount,
		CategoryID: transaction.CategoryID,
		// SubCategoryID: transaction.SubCategoryID,
		WalletID:    transaction.WalletID,
		Type:        transaction.Type,
		Description: transaction.Description,
		ImageUrl:    transaction.ImageUrl,
		Date:        transaction.Date,
		UserID:      transaction.UserID,
		Category:    nil,
		SubCategory: nil,
		Wallet:      nil,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}

	if transaction.CategoryID != 0 {
		formatCategory := FormatCreateCategoryResponse(transaction.Category)
		response.Category = formatCategory
	}

	// if transaction.SubCategoryID != 0 {
	// 	formatSubCategory := FormatGetSubCategoryResponse(transaction.SubCategory)
	// 	response.SubCategory = formatSubCategory
	// }

	if transaction.WalletID != 0 {
		formatWallet := FormatCreateWalletResponse(transaction.Wallet)
		response.Wallet = formatWallet
	}

	return response
}

func FormatGetTransactionsResponse(transactions []entity.Transaction) []GetTransactionResponse {
	var data []GetTransactionResponse

	for _, transaction := range transactions {
		format := FormatGetTransactionResponse(transaction)
		data = append(data, format)
	}

	return data
}
