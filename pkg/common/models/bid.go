package models

import (
	"gorm.io/gorm"
)

type Bid struct {
	gorm.Model
	Amount     float64 `json:"amount"`
	InvoiceId  int     `json:"invoice_id"`
	InvestorId int     `json:"investor_id"`
	Investor   User    `json:"investor"`
}
