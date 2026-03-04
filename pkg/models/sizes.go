package models

const TableNameSize = "sizes"

type Size struct {
	Width          int64    `gorm:"column:width;primaryKey"`
	Height         int64    `gorm:"column:height;primaryKey"`
	GateType       GateType `gorm:"type:gate_type;primaryKey"`
	WholesalePrice int64    `gorm:"column:wholesale_price"`
	RetailPrice    int64    `gorm:"column:retail_price"`
}

func (*Size) TableName() string {
	return TableNameSize
}
