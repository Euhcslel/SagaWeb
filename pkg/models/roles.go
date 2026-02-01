package models

const TableNameRoles = "roles"

type Role struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name string `gorm:"column:name;not null"`
}

func (*Role) TableName() string {
	return TableNameRoles
}
