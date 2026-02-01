package models

const TableNameSize = "sizes"

type Size struct {
	GateTypeID     int64 `gorm:"column:gate_type_id;primaryKey"`
	Width          int64 `gorm:"column:width;primaryKey"`
	Height         int64 `gorm:"column:height;primaryKey"`
	WholesalePrice int64 `gorm:"column:wholesale_price"`
	RetailPrice    int64 `gorm:"column:retail_price"`

	GateType GateType `gorm:"foreignKey:GateTypeID;references:ID"`
}

func (*Size) TableName() string {
	return TableNameSize
}
