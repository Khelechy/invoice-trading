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

func CreateUser(db *gorm.DB, user User) (uint, error) {
	result := db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

func GetIssuer(db *gorm.DB, issuerId string) (*User, error){
	var user User

	err := db.Model(&User{}).Where("id = ? AND user_type = ?", issuerId, "issuer").First(&user).Error
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
