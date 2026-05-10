package order_gate_residential_drives

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/rails"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
)

const TableNameOrderGateResidentialDrive = "order_gate_residential_drives"

type OrderGateResidentialDrive struct {
	OrderID   int64 `gorm:"column:order_id;primaryKey"`
	RowNumber int64 `gorm:"column:row_number;primaryKey"`
	DriveID   int64 `gorm:"column:drive_id;not null"`
	RailID    int64 `gorm:"column:rail_id;not null"`

	Drive residential_gate_drives.ResidentialGateDrive `gorm:"foreignKey:DriveID;references:ID"`
	Rail  rails.Rail                                   `gorm:"foreignKey:RailID;references:ID"`
}

func (*OrderGateResidentialDrive) TableName() string {
	return TableNameOrderGateResidentialDrive
}
