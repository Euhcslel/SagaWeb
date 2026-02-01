package models

import (
	"time"
)

const TableNameSalesAndOffer = "sales_and_offers"

type SalesAndOffer struct {
	SaleID      int64     `gorm:"column:sale_id;primaryKey"`
	OfferNumber string    `gorm:"column:offer_number;primaryKey"`
	Path        string    `gorm:"column:path;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (*SalesAndOffer) TableName() string {
	return TableNameSalesAndOffer
}
