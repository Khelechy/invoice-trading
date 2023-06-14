package models

type Invoice struct {
	Id          int     `json:"id" gorm:"primaryKey"`
	Amount      float64 `json:"amount"`
	IssuerId    string  `json:"issuer_id"`
	InvoiceDate string  `json:"invoice_date"`
	Reference   string  `json:"reference"`
	Status      string  `json:"status"`
	Bids        []Bid   `json:"bids"`
}

type CreateInvoiceDto struct {
	Amount   float64 `json:"amount"`
	IssuerId string  `json:"issuer_id"`
}
