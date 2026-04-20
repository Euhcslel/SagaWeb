package residential_gate_drives

import (
	"github.com/shopspring/decimal"
	"project/internal/domain/units"
)

const TableNameResidentialGateDrife = "residential_gate_drives"

type ResidentialGateDrive struct {
	ID             int64           `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name           string          `gorm:"column:name;not null"`
	UnitID         int64           `gorm:"column:unit_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	Specifications string          `gorm:"column:specifications"`

	Unit units.Unit `gorm:"foreignKey:UnitID;references:ID"`
}

func (*ResidentialGateDrive) TableName() string {
	return TableNameResidentialGateDrife
}
