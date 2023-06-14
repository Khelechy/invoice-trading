package models

type User struct {
	Id          int     `json:"id" gorm:"primaryKey"`
	Balance      float64 `json:"amount"`
	UserType    string  `json:"user_type"`
}
