package models

import (
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	Amount      int `json:"amount"`
	AmountBided int `json:"amount_bided"`
	IssuerId    uint     `json:"issuer_id"`
	Reference   string  `json:"reference"`
	Status      string  `json:"status"`
	Bids        []Bid   `gorm:"constraint:OnUpdate:CASCADE" json:"bids"`
}

type CreateInvoiceDto struct {
	Amount   int `json:"amount"`
	IssuerId uint     `json:"issuer_id"`
}
