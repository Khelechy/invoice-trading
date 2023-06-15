package models

import (
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	Amount      int    `json:"amount"`
	AmountBided int    `json:"amount_bided"`
	IssuerId    uint   `json:"issuer_id"`
	Reference   string `json:"reference"`
	Status      string `json:"status"`
	Bids        []Bid  `gorm:"constraint:OnUpdate:CASCADE" json:"bids"`
}

type CreateInvoiceDto struct {
	Amount   int  `json:"amount"`
	IssuerId uint `json:"issuer_id"`
}

func CreateInvoice(db *gorm.DB, invoice Invoice) (uint, error) {
	result := db.Create(&invoice)
	if result.Error != nil {
		return 0, result.Error
	}

	return invoice.ID, nil
}

func GetInvoice(db *gorm.DB, id string) (*Invoice, error) {
	var invoice Invoice

	err := db.Model(&Invoice{}).Where("id = ?", id).Preload("Bids.Investor").First(&invoice).Error
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}
