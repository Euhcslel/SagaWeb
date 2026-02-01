package models

const TableNameCompany = "companies"

type Company struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name string `gorm:"column:name;not null"`
}

func (*Company) TableName() string {
	return TableNameCompany
}
