package models

import (
	"time"
)

const TableNameMontagesAndSale = "montages_and_sales"

type MontagesAndSale struct {
	SaleID      int64     `gorm:"column:sale_id;primaryKey" json:"sale_id"`
	TypeID      int32     `gorm:"column:type_id" json:"type_id"`
	MontageDate time.Time `gorm:"column:montage_date;not null" json:"montage_date"`
	Price       int32     `gorm:"column:price" json:"price"`

	Type MontageType `gorm:"foreignKey:TypeID;references:ID"`
}

func (*MontagesAndSale) TableName() string {
	return TableNameMontagesAndSale
}
