package models

import (
	"time"
)

const TableNameSalesAndContract = "sales_and_contracts"

type SalesAndContract struct {
	SaleID         int64     `gorm:"column:sale_id;primaryKey" json:"sale_id"`
	ContractNumber string    `gorm:"column:contract_number;primaryKey" json:"contract_number"`
	Path           string    `gorm:"column:path;not null" json:"path"`
	CreatedAt      time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (*SalesAndContract) TableName() string {
	return TableNameSalesAndContract
}
