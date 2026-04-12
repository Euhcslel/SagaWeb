package products_price_history

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameProductsPriceHistory = "products_price_history"

type ProductsPriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	ProductID      int64           `gorm:"column:product_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*ProductsPriceHistory) TableName() string {
	return TableNameProductsPriceHistory
}
