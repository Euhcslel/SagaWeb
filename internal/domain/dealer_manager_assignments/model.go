// Package dealer_manager_assignments предоставляет модель для связи менеджера с дилером.
// Данные хранятся в таблице dealer_manager_assignments.

package dealer_manager_assignments

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/dealers"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
)

const TableNameDealerManagerAssignment = "dealer_manager_assignments"

type DealerManagerAssignment struct {
	ManagerID int64 `gorm:"column:manager_id;primaryKey"`
	DealerID  int64 `gorm:"column:dealer_id;primaryKey"`

	Manager users.User     `gorm:"foreignKey:ManagerID;references:ID"`
	Dealer  dealers.Dealer `gorm:"foreignKey:DealerID;references:UserID"`
}

func (*DealerManagerAssignment) TableName() string {
	return TableNameDealerManagerAssignment
}
