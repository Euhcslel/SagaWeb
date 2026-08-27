// Package rail_price_history предоставляет модель для работы с историей цен за направляющие.
// История цен за направляющие хранится в таблице rail_price_history.

package rail_price_history

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameRailPriceHistory = "rail_price_history"

type RailPriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	RailID         int64           `gorm:"column:rail_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*RailPriceHistory) TableName() string {
	return TableNameRailPriceHistory
}
