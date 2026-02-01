package models

const TableNameDealer = "dealers"

type Dealer struct {
	UserID    int64 `gorm:"column:user_id;primaryKey"`
	CompanyID int64 `gorm:"column:company_id"`

	User    User    `gorm:"foreignKey:UserID;references:ID"`
	Company Company `gorm:"foreignKey:CompanyID;references:ID"`
}

func (*Dealer) TableName() string {
	return TableNameDealer
}
