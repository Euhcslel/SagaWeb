package gate_sizes

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/shopspring/decimal"
)

const TableNameGateSize = "gate_sizes"

type GateSize struct {
	Width          int64           `gorm:"column:width;primaryKey"`
	Height         int64           `gorm:"column:height;primaryKey"`
	GateType       enums.GateType  `gorm:"column:gate_type;type:gate_type;primaryKey"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
}

func (*GateSize) TableName() string {
	return TableNameGateSize
}
