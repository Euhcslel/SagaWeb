package managers_and_dealers

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
)

const TableNameManagerAndDealer = "managers_and_dealers"

type ManagerAndDealer struct {
	ManagerID int64 `gorm:"column:manager_id;primaryKey"`
	DealerID  int64 `gorm:"column:dealer_id;primaryKey"`

	Manager users.User `gorm:"foreignKey:ManagerID;references:ID"`
	Dealer  users.User `gorm:"foreignKey:DealerID;references:ID"`
}

func (*ManagerAndDealer) TableName() string {
	return TableNameManagerAndDealer
}
