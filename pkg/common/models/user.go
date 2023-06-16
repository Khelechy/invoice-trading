package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Balance  int    `json:"balance"`
	UserType string `json:"user_type"`
}

func CreateUser(db *gorm.DB, user User) (*User, error) {
	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func GetIssuer(db *gorm.DB, issuerId uint) (*User, error) {
	var user User

	err := db.Model(&User{}).Where("id = ? AND user_type = ?", issuerId, "issuer").First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetInvestor(db *gorm.DB, investorId uint) (*User, error) {
	var user User

	err := db.Model(&User{}).Where("id = ? AND user_type = ?", investorId, "investor").First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetInvestors(db *gorm.DB) ([]User, error) {
	var users []User

	err := db.Model(&User{}).Where("user_type = ?", "investor").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
