package models

const TableNameIndustrialGatesAndSalesDrive = "industrial_gates_and_sales_drive"

type IndustrialGatesAndSalesDrive struct {
	SaleID    int64 `gorm:"column:sale_id;primaryKey" json:"sale_id"`
	RowNumber int64 `gorm:"column:row_number;primaryKey" json:"row_number"`
	DriveID   int32 `gorm:"column:drive_id;not null" json:"drive_id"`

	IndustrialGateDrive IndustrialGateDrive `gorm:"foreignKey:DriveID;references:ID"`
}

func (*IndustrialGatesAndSalesDrive) TableName() string {
	return TableNameIndustrialGatesAndSalesDrive
}
