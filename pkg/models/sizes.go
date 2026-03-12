package models

import "github.com/shopspring/decimal"

const TableNameSize = "sizes"

type Size struct {
	Width          int64           `gorm:"column:width;primaryKey"`
	Height         int64           `gorm:"column:height;primaryKey"`
	GateType       GateType        `gorm:"type:gate_type;primaryKey"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
}

func (*Size) TableName() string {
	return TableNameSize
}
