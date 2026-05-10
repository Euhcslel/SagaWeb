package order_gate_industrial_drives

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
)

const TableNameOrderGateIndustrialDrive = "order_gate_industrial_drives"

type OrderGateIndustrialDrive struct {
	OrderID   int64 `gorm:"column:order_id;primaryKey"`
	RowNumber int64 `gorm:"column:row_number;primaryKey"`
	DriveID   int64 `gorm:"column:drive_id;not null"`

	IndustrialGateDrive industrial_gate_drives.IndustrialGateDrive `gorm:"foreignKey:DriveID;references:ID"`
}

func (*OrderGateIndustrialDrive) TableName() string {
	return TableNameOrderGateIndustrialDrive
}
