package models

const TableNameRoles = "roles"

type Role struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`
}

func (*Role) TableName() string {
	return TableNameRoles
}
