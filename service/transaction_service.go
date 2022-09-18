package service

import (
	"errors"
	"pawang-backend/entity"
	"pawang-backend/model/request"
	"pawang-backend/repository"
)

type TransactionService interface {
	CreateTransaction(userID int, input request.CreateTransactionRequest) (entity.Transaction, error)
	GetTransactions(userID int) ([]entity.Transaction, error)
	GetTransactionByID(transactionID int, userID int) (entity.Transaction, error)
	UpdateTransaction(transactionID int, userID int, input request.CreateTransactionRequest) (entity.Transaction, error)
	DeleteTransaction(transactionID int, userID int) error
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	walletRepository      repository.WalletRepository
	subCategoryRepository repository.SubCategoryRepository
	categoryRepository    repository.CategoryRepository
}

func NewTransactionService(transactionRepository repository.TransactionRepository, walletRepository repository.WalletRepository, subCategoryRepository repository.SubCategoryRepository, categoryRepository repository.CategoryRepository) *transactionService {
	return &transactionService{transactionRepository, walletRepository, subCategoryRepository, categoryRepository}
}

func (service *transactionService) CreateTransaction(userID int, input request.CreateTransactionRequest) (entity.Transaction, error) {
	transaction := entity.Transaction{}
	transaction.Amount = input.Amount
	transaction.CategoryID = input.CategoryID
	transaction.WalletID = input.WalletID
	transaction.Type = input.Type
	transaction.Description = input.Description
	transaction.ImageUrl = input.ImageUrl
	transaction.Date = input.Date
	transaction.UserID = userID

	// Check Type Transaction same Type Category
	category, err := service.categoryRepository.FindByID(userID, input.CategoryID)
	if err != nil {
		return transaction, err
	}

	if category.Type != input.Type {
		return transaction, errors.New("category type mismatch, please check your type transaction")
	}

	// if input.SubCategoryID != 0 {
	// 	// Check Sub Category ID on Main Category
	// 	subCategory, err := service.subCategoryRepository.FindByID(input.SubCategoryID)
	// 	if err != nil {
	// 		return transaction, err
	// 	}

	// 	if subCategory.CategoryID != input.CategoryID {
	// 		return transaction, errors.New("sub category does not exist on main category")
	// 	}

	// 	transaction.SubCategoryID = input.SubCategoryID
	// }

	wallet, err := service.walletRepository.FindByID(input.WalletID)
	if err != nil {
		return transaction, err
	}

	if wallet.ID == 0 {
		return transaction, errors.New("the wallet does not exist")
	}

	if wallet.UserID != userID {
		return transaction, errors.New("the wallet does not exist")
	}

	if transaction.Type == "income" {
		wallet.Balance = wallet.Balance + transaction.Amount
		_, err := service.walletRepository.Update(wallet)
		if err != nil {
			return transaction, err
		}
	} else if transaction.Type == "outcome" {
		wallet.Balance = wallet.Balance - transaction.Amount
		_, err := service.walletRepository.Update(wallet)
		if err != nil {
			return transaction, err
		}
	} else {
		return transaction, errors.New("type does not exist")
	}

	newTransaction, err := service.transactionRepository.Insert(transaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (service *transactionService) GetTransactions(userID int) ([]entity.Transaction, error) {
	transactions, err := service.transactionRepository.FindAllByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (service *transactionService) GetTransactionByID(transactionID int, userID int) (entity.Transaction, error) {
	transaction, err := service.transactionRepository.FindByID(transactionID)
	if err != nil {
		return transaction, err
	}

	if transaction.ID == 0 {
		return transaction, errors.New("the transaction does not exist")
	}

	if transaction.UserID != userID {
		return transaction, errors.New("the transaction does not exist")
	}

	return transaction, nil
}

func (service *transactionService) UpdateTransaction(transactionID int, userID int, input request.CreateTransactionRequest) (entity.Transaction, error) {
	transaction, err := service.transactionRepository.FindByID(transactionID)
	if err != nil {
		return transaction, err
	}

	if transaction.ID == 0 {
		return transaction, errors.New("the transaction does not exist")
	}

	if transaction.UserID != userID {
		return transaction, errors.New("the transaction does not exist")
	}

	// Check Type Transaction same Type Category
	category, err := service.categoryRepository.FindByID(userID, input.CategoryID)
	if err != nil {
		return transaction, err
	}

	if category.Type != input.Type {
		return transaction, errors.New("category type mismatch, please check your type transaction")
	}

	wallet, err := service.walletRepository.FindByID(input.WalletID)
	if err != nil {
		return transaction, err
	}

	if wallet.ID == 0 {
		return transaction, errors.New("the wallet does not exist")
	}

	if wallet.UserID != userID {
		return transaction, errors.New("the wallet does not exist")
	}

	if transaction.Type == "income" {
		wallet.Balance = (wallet.Balance - transaction.Amount) + input.Amount
		_, err := service.walletRepository.Update(wallet)
		if err != nil {
			return transaction, err
		}
	} else if transaction.Type == "outcome" {
		wallet.Balance = (wallet.Balance + transaction.Amount) - input.Amount
		_, err := service.walletRepository.Update(wallet)
		if err != nil {
			return transaction, err
		}
	} else {
		return transaction, errors.New("type does not exist")
	}

	transaction.Amount = input.Amount
	transaction.CategoryID = input.CategoryID
	transaction.WalletID = input.WalletID
	transaction.Type = input.Type
	transaction.Description = input.Description
	transaction.ImageUrl = input.ImageUrl
	transaction.Date = input.Date
	transaction.UserID = userID

	// if input.SubCategoryID != 0 {
	// 	subCategory, err := service.subCategoryRepository.FindByID(input.SubCategoryID)
	// 	if err != nil {
	// 		return transaction, err
	// 	}

	// 	if subCategory.CategoryID != input.CategoryID {
	// 		return transaction, errors.New("sub category does not exist on main category")
	// 	}

	// 	transaction.SubCategoryID = input.SubCategoryID
	// }

	updateTransaction, err := service.transactionRepository.Update(transaction)
	if err != nil {
		return updateTransaction, err
	}

	return updateTransaction, nil
}

func (service *transactionService) DeleteTransaction(transactionID int, userID int) error {
	transaction, err := service.transactionRepository.FindByID(transactionID)
	if err != nil {
		return err
	}

	if transaction.ID == 0 {
		return errors.New("the transaction does not exist")
	}

	if transaction.UserID != userID {
		return errors.New("the transaction does not exist")
	}

	wallet, err := service.walletRepository.FindByID(transaction.WalletID)
	if err != nil {
		return err
	}

	if wallet.ID == 0 {
		return errors.New("the wallet does not exist")
	}

	if wallet.UserID != userID {
		return errors.New("the wallet does not exist")
	}

	if transaction.Type == "income" {
		wallet.Balance = (wallet.Balance - transaction.Amount)
		_, err := service.walletRepository.Update(wallet)
		if err != nil {
			return err
		}
	} else if transaction.Type == "outcome" {
		wallet.Balance = (wallet.Balance + transaction.Amount)
		_, err := service.walletRepository.Update(wallet)
		if err != nil {
			return err
		}
	} else {
		return errors.New("type does not exist")
	}

	err = service.transactionRepository.Delete(transaction)
	if err != nil {
		return err
	}

	return nil
}
