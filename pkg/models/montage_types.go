package models

const TableNameMontageType = "montage_types"

type MontageType struct {
	ID   int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`
}

func (*MontageType) TableName() string {
	return TableNameMontageType
}
