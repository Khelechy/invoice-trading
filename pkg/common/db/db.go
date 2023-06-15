package db

import (
	"log"

	"github.com/khelechy/invoice-trading/pkg/common/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.Debug().AutoMigrate(&models.Invoice{}, &models.Bid{}, &models.User{})

	return db
}
