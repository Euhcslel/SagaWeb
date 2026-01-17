package models

const TableNameSize = "sizes"

type Size struct {
	GateTypeID int64 `gorm:"column:gate_type_id;primaryKey" json:"gate_type_id"`
	Width      int64 `gorm:"column:width;primaryKey" json:"width"`
	Height     int64 `gorm:"column:height;primaryKey" json:"height"`
	WholesalePrice      int64 `gorm:"column:wholesale_price" json:"wholesale_price"`
	RetailPrice int64 `gorm:"column:retail_price" json:"retail_price"`

	GateType GateType `gorm:"foreignKey:GateTypeID;references:ID"`
}

func (*Size) TableName() string {
	return TableNameSize
}
