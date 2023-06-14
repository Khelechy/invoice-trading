package models

type Bid struct {
	Id         int     `json:"id" gorm:"primaryKey"`
	Amount     float64 `json:"amount"`
	InvestorId string  `json:"investor_id"`
	Investor User `json:"investor"`
}
