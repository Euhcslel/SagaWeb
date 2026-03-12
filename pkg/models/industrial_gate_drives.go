package models

import "github.com/shopspring/decimal"

const TableNameIndustrialGateDrife = "industrial_gate_drives"

type IndustrialGateDrive struct {
	ID             int32           `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name           string          `gorm:"column:name;not null"`
	UnitID         int32           `gorm:"column:unit_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	Specifications string          `gorm:"column:specifications"`

	Unit Unit `gorm:"foreignKey:UnitID;references:ID"`
}

func (*IndustrialGateDrive) TableName() string {
	return TableNameIndustrialGateDrife
}
