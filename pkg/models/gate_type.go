package models

const TableNameGateType = "gate_type"

type GateType struct {
	ID   int32  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name string `gorm:"column:name;not null"`
}

func (*GateType) TableName() string {
	return TableNameGateType
}
