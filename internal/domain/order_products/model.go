// Package order_products предоставляет модель для работы с товарами заказа.
// Данные о товарах заказа хранятся в таблице order_products.

package order_products

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
)

const TableNameOrderProduct = "order_products"

type OrderProduct struct {
	OrderID   int64 `gorm:"column:order_id;primaryKey"`
	ProductID int64 `gorm:"column:product_id;primaryKey"`
	Amount    int32 `gorm:"column:amount;not null"`

	Product products.Product `gorm:"foreignKey:ProductID;references:ID"`
}

func (*OrderProduct) TableName() string {
	return TableNameOrderProduct
}
