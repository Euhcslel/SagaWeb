package models

const TableNameSalesAndProduct = "sales_and_products"

type SalesAndProduct struct {
	SaleID    int64 `gorm:"column:sale_id;primaryKey"`
	ProductID int32 `gorm:"column:product_id;primaryKey"`
	Amount    int32 `gorm:"column:amount;not null"`

	Product Product `gorm:"foreignKey:ProductID;references:ID"`
}

func (*SalesAndProduct) TableName() string {
	return TableNameSalesAndProduct
}
