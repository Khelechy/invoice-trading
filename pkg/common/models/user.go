package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Balance      float64 `json:"amount"`
	UserType    string  `json:"user_type"`
}
