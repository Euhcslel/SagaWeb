package sales_and_contracts

import (
	"time"
)

const TableNameSalesAndContract = "sales_and_contracts"

type SalesAndContract struct {
	SaleID         int64     `gorm:"column:sale_id;primaryKey"`
	ContractNumber string    `gorm:"column:contract_number;primaryKey"`
	Path           string    `gorm:"column:path;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (*SalesAndContract) TableName() string {
	return TableNameSalesAndContract
}
