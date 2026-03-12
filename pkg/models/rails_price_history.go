package models

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameRailsPriceHistory = "rails_price_history"

type RailsPriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	RailID         int64           `gorm:"column:rail_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*RailsPriceHistory) TableName() string {
	return TableNameRailsPriceHistory
}
