package order_gate_options

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
)

const TableNameOrderGateOption = "order_gate_options"

type OrderGateOption struct {
	OrderID   int64 `gorm:"column:order_id;primaryKey"`
	RowNumber int64 `gorm:"column:row_number;primaryKey"`
	OptionID  int64 `gorm:"column:option_id;primaryKey"`
	Amount    int32 `gorm:"column:amount;not null"`

	Option options.Option `gorm:"foreignKey:OptionID;references:ID"`
}

func (*OrderGateOption) TableName() string {
	return TableNameOrderGateOption
}
