package models

const TableNameUnit = "units"

type Unit struct {
	ID   int32  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name string `gorm:"column:name;not null"`
}

func (*Unit) TableName() string {
	return TableNameUnit
}
