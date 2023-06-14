package models

import (
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	Amount      float64 `json:"amount"`
	AmountBided float64 `json:"amount_bided"`
	IssuerId    int     `json:"issuer_id"`
	InvoiceDate string  `json:"invoice_date"`
	Reference   string  `json:"reference"`
	Status      string  `json:"status"`
	Bids        []Bid   `json:"bids"`
}

type CreateInvoiceDto struct {
	Amount   float64 `json:"amount"`
	IssuerId int     `json:"issuer_id"`
}
