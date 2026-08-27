// Package order_offers предоставляет модель для работы с коммерческими предложениями.
// Данные о каждом коммерческом предложении хранятся в таблице order_offers.

package order_offers

import (
	"time"
)

const TableNameOrderOffer = "order_offers"

type OrderOffer struct {
	OrderID     int64     `gorm:"column:order_id;primaryKey"`
	OfferNumber string    `gorm:"column:offer_number;primaryKey"`
	Path        string    `gorm:"column:path;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (*OrderOffer) TableName() string {
	return TableNameOrderOffer
}
