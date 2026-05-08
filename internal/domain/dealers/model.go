package dealers

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/companies"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
)

const TableNameDealer = "dealers"

type Dealer struct {
	UserID    int64  `gorm:"column:user_id;primaryKey"`
	CompanyID int64  `gorm:"column:company_id"`
	Address   string `gorm:"column:address"`

	User    users.User        `gorm:"foreignKey:UserID;references:ID"`
	Company companies.Company `gorm:"foreignKey:CompanyID;references:ID"`
}

func (*Dealer) TableName() string {
	return TableNameDealer
}
