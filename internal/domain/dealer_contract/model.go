// Package dealer_contract предоставляет модель для работы с договорами дилеров.
// Данные о договоре каждого дилера хранятся в таблице dealer_contract.

package dealer_contract

import (
	"time"
)

const TableNameDealerContract = "dealer_contract"

type DealerContract struct {
	DealerID       int64     `gorm:"column:dealer_id;primaryKey"`
	ContractNumber string    `gorm:"column:contract_number;not null"`
	SignedAt       time.Time `gorm:"column:signed_at;not null"`
	Path           *string   `gorm:"column:path"`
	CreatedAt      time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (*DealerContract) TableName() string {
	return TableNameDealerContract
}
