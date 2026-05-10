package product_price_history

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameProductPriceHistory = "product_price_history"

type ProductPriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	ProductID      int64           `gorm:"column:product_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*ProductPriceHistory) TableName() string {
	return TableNameProductPriceHistory
}
