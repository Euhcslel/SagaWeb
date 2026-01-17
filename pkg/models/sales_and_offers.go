package models

import (
	"time"
)

const TableNameSalesAndOffer = "sales_and_offers"

type SalesAndOffer struct {
	SaleID      int64     `gorm:"column:sale_id;primaryKey" json:"sale_id"`
	OfferNumber string    `gorm:"column:offer_number;primaryKey" json:"offer_number"`
	Path        string    `gorm:"column:path;not null" json:"path"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (*SalesAndOffer) TableName() string {
	return TableNameSalesAndOffer
}
