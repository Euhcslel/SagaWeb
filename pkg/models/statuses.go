package models

const TableNameStatus = "statuses"

type Status struct {
	ID   int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`
}

func (*Status) TableName() string {
	return TableNameStatus
}
