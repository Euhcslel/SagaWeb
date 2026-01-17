package models

const TableNameGateType = "gate_type"

type GateType struct {
	ID   int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`
}

func (*GateType) TableName() string {
	return TableNameGateType
}
