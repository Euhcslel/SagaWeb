package industrial_gates_and_sales_drive

import (
	"project/internal/domain/industrial_gate_drives"
)

const TableNameIndustrialGatesAndSalesDrive = "industrial_gates_and_sales_drive"

type IndustrialGatesAndSalesDrive struct {
	SaleID    int64 `gorm:"column:sale_id;primaryKey"`
	RowNumber int64 `gorm:"column:row_number;primaryKey"`
	DriveID   int32 `gorm:"column:drive_id;not null"`

	IndustrialGateDrive industrial_gate_drives.IndustrialGateDrive `gorm:"foreignKey:DriveID;references:ID"`
}

func (*IndustrialGatesAndSalesDrive) TableName() string {
	return TableNameIndustrialGatesAndSalesDrive
}
