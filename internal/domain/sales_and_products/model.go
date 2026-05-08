package sales_and_products

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
)

const TableNameSalesAndProduct = "sales_and_products"

type SalesAndProduct struct {
	SaleID    int64 `gorm:"column:sale_id;primaryKey"`
	ProductID int64 `gorm:"column:product_id;primaryKey"`
	Amount    int32 `gorm:"column:amount;not null"`

	Product products.Product `gorm:"foreignKey:ProductID;references:ID"`
}

func (*SalesAndProduct) TableName() string {
	return TableNameSalesAndProduct
}
