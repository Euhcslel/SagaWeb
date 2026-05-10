package gate_size_price_history

import (
	"time"

	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/shopspring/decimal"
)

const TableNameGateSizePriceHistory = "gate_size_price_history"

type GateSizePriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	Width          int64           `gorm:"column:width"`
	Height         int64           `gorm:"column:height"`
	GateType       enums.GateType  `gorm:"column:gate_type;type:gate_type"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*GateSizePriceHistory) TableName() string {
	return TableNameGateSizePriceHistory
}
