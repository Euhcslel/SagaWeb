package models

const TableNameGatesAndSalesOption = "gates_and_sales_options"

type GatesAndSalesOption struct {
	SaleID    int64 `gorm:"column:sale_id;primaryKey" json:"sale_id"`
	RowNumber int64 `gorm:"column:row_number;primaryKey" json:"row_number"`
	OptionID  int32 `gorm:"column:option_id;not null" json:"option_id"`
	Amount    int32 `gorm:"column:amount;not null" json:"amount"`

	Option Option `gorm:"foreignKey:OptionID;references:ID"`
}

func (*GatesAndSalesOption) TableName() string {
	return TableNameGatesAndSalesOption
}
