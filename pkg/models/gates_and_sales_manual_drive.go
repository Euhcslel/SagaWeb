package models

const TableNameGatesAndSalesManualDrive = "gates_and_sales_manual_drive"

type GatesAndSalesManualDrive struct {
	SaleID      int64 `gorm:"column:sale_id;primaryKey" json:"sale_id"`
	RowNumber   int64 `gorm:"column:row_number;primaryKey" json:"row_number"`
	ChainLength int32 `gorm:"column:chain_length" json:"chain_length"`
}

func (*GatesAndSalesManualDrive) TableName() string {
	return TableNameGatesAndSalesManualDrive
}
