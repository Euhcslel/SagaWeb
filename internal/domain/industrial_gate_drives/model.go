package industrial_gate_drives

import (
	"github.com/shopspring/decimal"
)

const TableNameIndustrialGateDrive = "industrial_gate_drives"

type IndustrialGateDrive struct {
	ID             int64           `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name           string          `gorm:"column:name;not null"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	Specifications *string         `gorm:"column:specifications"`
}

func (*IndustrialGateDrive) TableName() string {
	return TableNameIndustrialGateDrive
}
