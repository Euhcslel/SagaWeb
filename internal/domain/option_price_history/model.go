// Package option_price_history предоставляет модель для работы с историей цен за дополнительные опции.
// История цен за дополнительные опции хранится в таблице option_price_history.

package option_price_history

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameOptionPriceHistory = "option_price_history"

type OptionPriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	OptionID       int64           `gorm:"column:option_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*OptionPriceHistory) TableName() string {
	return TableNameOptionPriceHistory
}
