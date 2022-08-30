package service

import (
	"errors"
	"pawang-backend/entity"
	"pawang-backend/model/request"
	"pawang-backend/repository"
	"time"
)

type WalletService interface {
	GetWallets(userID int) ([]entity.Wallet, error)
	CreateWallet(userID int, input request.CreateWalletRequest) (entity.Wallet, error)
	UpdateWallet(walletID int, userID int, input request.UpdateWalletRequest) (entity.Wallet, error)
	DeleteWallet(walletID int, userID int) error
}

type walletService struct {
	walletRepository      repository.WalletRepository
	transactionRepository repository.TransactionRepository
}

func NewWalletService(walletRepository repository.WalletRepository, transactionRepository repository.TransactionRepository) *walletService {
	return &walletService{walletRepository, transactionRepository}
}

func (service *walletService) GetWallets(userID int) ([]entity.Wallet, error) {
	wallets, err := service.walletRepository.FindAllByUserID(userID)
	if err != nil {
		return wallets, err
	}

	return wallets, nil
}

func (service *walletService) CreateWallet(userID int, input request.CreateWalletRequest) (entity.Wallet, error) {
	wallet := entity.Wallet{}
	wallet.UserID = userID
	wallet.Name = input.Name
	wallet.Balance = input.Balance

	newWallet, err := service.walletRepository.Insert(wallet)
	if err != nil {
		return wallet, err
	}

	if input.Balance > 0 {
		transaction := entity.Transaction{}
		transaction.Amount = input.Balance
		transaction.CategoryID = 12
		transaction.Date = time.Now()
		transaction.WalletID = newWallet.ID
		transaction.Type = "income"
		transaction.UserID = userID

		_, err := service.transactionRepository.Insert(transaction)
		if err != nil {
			return wallet, err
		}
	}

	return newWallet, err
}

func (service *walletService) UpdateWallet(walletID int, userID int, input request.UpdateWalletRequest) (entity.Wallet, error) {
	wallet, err := service.walletRepository.FindByID(walletID)
	if err != nil {
		return wallet, err
	}

	if wallet.ID == 0 {
		return wallet, errors.New("the wallet does not exist")
	}

	if wallet.UserID != userID {
		return wallet, errors.New("the wallet does not exist")
	}

	if wallet.Balance < input.Balance {
		transaction := entity.Transaction{}
		transaction.Amount = input.Balance - wallet.Balance
		transaction.CategoryID = 12
		transaction.Date = time.Now()
		transaction.WalletID = walletID
		transaction.Type = "income"
		transaction.UserID = userID

		_, err := service.transactionRepository.Insert(transaction)
		if err != nil {
			return wallet, err
		}
	} else if wallet.Balance > input.Balance {
		transaction := entity.Transaction{}
		transaction.Amount = wallet.Balance - input.Balance
		transaction.CategoryID = 13
		transaction.Date = time.Now()
		transaction.WalletID = walletID
		transaction.Type = "outcome"
		transaction.UserID = userID

		_, err := service.transactionRepository.Insert(transaction)
		if err != nil {
			return wallet, err
		}
	}

	wallet.Name = input.Name
	wallet.Balance = input.Balance

	updateWallet, err := service.walletRepository.Update(wallet)
	if err != nil {
		return updateWallet, err
	}

	return updateWallet, nil
}

func (service *walletService) DeleteWallet(walletID int, userID int) error {
	wallet, err := service.walletRepository.FindByID(walletID)
	if err != nil {
		return err
	}

	if wallet.ID == 0 {
		return errors.New("the wallet does not exist")
	}

	if wallet.UserID != userID {
		return errors.New("the wallet does not exist")
	}

	err = service.walletRepository.Delete(wallet)
	if err != nil {
		return err
	}

	return nil
}
