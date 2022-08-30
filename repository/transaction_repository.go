package repository

import (
	"pawang-backend/entity"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Insert(transaction entity.Transaction) (entity.Transaction, error)
	FindAllByUserID(userID int) ([]entity.Transaction, error)
	FindByID(transactionID int) (entity.Transaction, error)
	Update(transaction entity.Transaction) (entity.Transaction, error)
	Delete(transaction entity.Transaction) error
}

type transactionRepository struct {
	database *gorm.DB
}

func NewTransactionRepository(database *gorm.DB) *transactionRepository {
	return &transactionRepository{database}
}

func (repository *transactionRepository) Insert(transaction entity.Transaction) (entity.Transaction, error) {
	if err := repository.database.Create(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (repository *transactionRepository) FindAllByUserID(userID int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := repository.database.Where("user_id = ?", userID).Order("created_at desc").Preload("Category").Preload("SubCategory").Preload("Wallet").Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (repository *transactionRepository) FindByID(transactionID int) (entity.Transaction, error) {
	var transaction entity.Transaction

	if err := repository.database.Where("id = ?", transactionID).Find(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (repository *transactionRepository) Update(transaction entity.Transaction) (entity.Transaction, error) {
	if err := repository.database.Save(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (repository *transactionRepository) Delete(transaction entity.Transaction) error {
	if err := repository.database.Delete(&transaction).Error; err != nil {
		return err
	}

	return nil
}
