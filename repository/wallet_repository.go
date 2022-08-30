package repository

import (
	"pawang-backend/entity"

	"gorm.io/gorm"
)

type WalletRepository interface {
	FindAllByUserID(userID int) ([]entity.Wallet, error)
	FindByID(walletID int) (entity.Wallet, error)
	Insert(wallet entity.Wallet) (entity.Wallet, error)
	Update(wallet entity.Wallet) (entity.Wallet, error)
	Delete(wallet entity.Wallet) error
}

type walletRepository struct {
	database *gorm.DB
}

func NewWalletRepository(database *gorm.DB) *walletRepository {
	return &walletRepository{database}
}

func (repository *walletRepository) Insert(wallet entity.Wallet) (entity.Wallet, error) {
	if err := repository.database.Create(&wallet).Error; err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (repository *walletRepository) FindAllByUserID(userID int) ([]entity.Wallet, error) {
	var wallets []entity.Wallet

	if err := repository.database.Where("user_id = ?", userID).Preload("Transactions").Find(&wallets).Error; err != nil {
		return wallets, err
	}

	return wallets, nil
}

func (repository *walletRepository) FindByID(walletID int) (entity.Wallet, error) {
	var wallet entity.Wallet

	if err := repository.database.Where("id = ?", walletID).Find(&wallet).Error; err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (repository *walletRepository) Update(wallet entity.Wallet) (entity.Wallet, error) {
	if err := repository.database.Save(&wallet).Error; err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (repository *walletRepository) Delete(wallet entity.Wallet) error {
	if err := repository.database.Delete(&wallet).Error; err != nil {
		return err
	}

	return nil
}
