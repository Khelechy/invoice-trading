package models

import (
	"gorm.io/gorm"
)

type Bid struct {
	gorm.Model
	Amount     int `json:"amount"`
	InvoiceId  uint     `json:"invoice_id"`
	InvestorId uint     `json:"investor_id"`
	Investor   User    `gorm:"foreignKey:InvestorId;" json:"investor"`
}

type PlaceBid struct {
	Amount     int `json:"amount"`
	InvoiceId  uint     `json:"invoice_id"`
	InvestorId uint     `json:"investor_id"`
}
