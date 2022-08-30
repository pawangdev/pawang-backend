package response

import (
	"pawang-backend/entity"
	"time"
)

type CreateWalletResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FormatCreateWalletResponse(wallet entity.Wallet) CreateWalletResponse {
	response := CreateWalletResponse{
		ID:        wallet.ID,
		Name:      wallet.Name,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}

	return response
}

type GetWalletResponse struct {
	ID           int                         `json:"id"`
	Name         string                      `json:"name"`
	UserID       int                         `json:"user_id"`
	Balance      int                         `json:"balance"`
	Transactions []CreateTransactionResponse `json:"transactions"`
	CreatedAt    time.Time                   `json:"created_at"`
	UpdatedAt    time.Time                   `json:"updated_at"`
}

func FormatGetWalletResponse(wallet entity.Wallet) GetWalletResponse {
	response := GetWalletResponse{
		ID:        wallet.ID,
		Name:      wallet.Name,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}

	var data []CreateTransactionResponse

	for _, transaction := range wallet.Transactions {
		formatter := FormatCreateTransactionResponse(transaction)
		data = append(data, formatter)
	}

	response.Transactions = data

	return response
}

func FormatGetWalletsResponse(wallets []entity.Wallet) []GetWalletResponse {
	walletsResponse := []GetWalletResponse{}

	for _, wallet := range wallets {
		walletFormat := FormatGetWalletResponse(wallet)
		walletsResponse = append(walletsResponse, walletFormat)
	}

	return walletsResponse
}
