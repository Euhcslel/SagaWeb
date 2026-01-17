package models

const TableNameUnit = "units"

type Unit struct {
	ID   int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`
}

func (*Unit) TableName() string {
	return TableNameUnit
}
