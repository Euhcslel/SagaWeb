package sizes_price_history

import (
	"time"

	"github.com/shopspring/decimal"
	"project/internal/domain/enums"
)

const TableNameSizesPriceHistory = "sizes_price_history"

type SizesPriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	Width          int64           `gorm:"column:width"`
	Height         int64           `gorm:"column:height"`
	GateType       enums.GateType  `gorm:"column:gate_type;type:gate_type"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*SizesPriceHistory) TableName() string {
	return TableNameSizesPriceHistory
}
