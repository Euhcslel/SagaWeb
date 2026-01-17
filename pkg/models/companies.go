package models

const TableNameCompany = "companies"

type Company struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`
}

func (*Company) TableName() string {
	return TableNameCompany
}
