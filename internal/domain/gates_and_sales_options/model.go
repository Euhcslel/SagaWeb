package gates_and_sales_options

import (
	"project/internal/domain/options"
)

const TableNameGatesAndSalesOption = "gates_and_sales_options"

type GatesAndSalesOption struct {
	SaleID    int64 `gorm:"column:sale_id;primaryKey"`
	RowNumber int64 `gorm:"column:row_number;primaryKey"`
	OptionID  int32 `gorm:"column:option_id;primaryKey"`
	Amount    int32 `gorm:"column:amount;not null"`

	Option options.Option `gorm:"foreignKey:OptionID;references:ID"`
}

func (*GatesAndSalesOption) TableName() string {
	return TableNameGatesAndSalesOption
}
