package models

const TableNameStatus = "statuses"

type Status struct {
	ID   int32  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name string `gorm:"column:name;not null"`
}

func (*Status) TableName() string {
	return TableNameStatus
}
