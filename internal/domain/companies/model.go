// Package companies предоставляет модель для работы с компаниями.
// Компании хранятся в таблице companies.

package companies

const TableNameCompany = "companies"

type Company struct {
	ID   int64   `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name string  `gorm:"column:name;not null"`
	INN  *string `gorm:"column:inn"`
	KPP  *string `gorm:"column:kpp"`
	OGRN *string `gorm:"column:ogrn"`
}

func (*Company) TableName() string {
	return TableNameCompany
}
