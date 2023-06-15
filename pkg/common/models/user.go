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
