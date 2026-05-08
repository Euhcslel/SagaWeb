package sizes

import (
	"github.com/shopspring/decimal"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
)

const TableNameSize = "sizes"

type Size struct {
	Width          int64           `gorm:"column:width;primaryKey"`
	Height         int64           `gorm:"column:height;primaryKey"`
	GateType       enums.GateType  `gorm:"type:gate_type;primaryKey"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
}

func (*Size) TableName() string {
	return TableNameSize
}
