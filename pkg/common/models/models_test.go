package models

import (
	"log"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	setUpDatabase()
	os.Exit(m.Run())
}

func setUpDatabase() {

	psqldb, err := gorm.Open(postgres.Open("postgres://postgres:root@localhost:5432/invoice_trading_db_test"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db = psqldb
}

func refreshUserTable() {
	db.Migrator().DropTable(&User{})
	db.Migrator().CreateTable(&User{})
}

func seedInvestorandIssuer() (*User, *User, error) {
	investor := User{
		Name:     "Kelechi",
		Balance:  5000,
		UserType: "investor",
	}

	issuer := User{
		Name:     "KCorp",
		Balance:  0,
		UserType: "issuer",
	}

	err := db.Model(&User{}).Create(&investor).Error

	if err != nil {
		log.Fatal("cannot seed investor")
	}

	err = db.Model(&User{}).Create(&issuer).Error

	if err != nil {
		log.Fatal("cannot seed issuer")
	}

	return &investor, &issuer, nil
}

func seedInvestors() (*User, *User, error) {
	investor1 := User{
		Name:     "John",
		Balance:  5000,
		UserType: "investor",
	}

	investor2 := User{
		Name:     "Doe",
		Balance:  0,
		UserType: "investor",
	}

	err := db.Model(&User{}).Create(&investor1).Error

	if err != nil {
		log.Fatal("cannot seed investor")
	}

	err = db.Model(&User{}).Create(&investor2).Error

	if err != nil {
		log.Fatal("cannot seed investor")
	}

	return &investor1, &investor2, nil
}
