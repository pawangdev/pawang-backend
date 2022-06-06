package config

import (
	"log"
	"pawang-backend/models"

	"gorm.io/gorm"
)

var transactionType = []models.TransactionType{
	models.TransactionType{
		Name: "income",
	},
	models.TransactionType{
		Name: "outcome",
	},
}

func Load(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.TransactionType{})

	if err != nil {
		log.Fatalf("Cannot Drop Table %v", err)
	}

	err = db.Debug().AutoMigrate(&models.TransactionType{})
	if err != nil {
		log.Fatalf("Cannot Migrate Table %v", err)
	}

	for i, _ := range transactionType {
		err = db.Debug().Model(&models.TransactionType{}).Create(&transactionType[i]).Error
		if err != nil {
			log.Fatalf("cannot seed transaction type %v", err)
		}
	}

}
