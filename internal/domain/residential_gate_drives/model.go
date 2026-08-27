// Package residential_gate_drives предоставляет модель для работы с бытовыми приводами.
// Данные о бытовых приводах хранятся в таблице residential_gate_drives.

package residential_gate_drives

import (
	"github.com/shopspring/decimal"
)

const TableNameResidentialGateDrive = "residential_gate_drives"

type ResidentialGateDrive struct {
	ID             int64           `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name           string          `gorm:"column:name;not null"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	Specifications *string         `gorm:"column:specifications"`
	ImagePath      string          `gorm:"column:image_path"`
}

func (*ResidentialGateDrive) TableName() string {
	return TableNameResidentialGateDrive
}
