package options_price_history

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameOptionsPriceHistory = "options_price_history"

type OptionsPriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	OptionID       int64           `gorm:"column:option_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*OptionsPriceHistory) TableName() string {
	return TableNameOptionsPriceHistory
}
