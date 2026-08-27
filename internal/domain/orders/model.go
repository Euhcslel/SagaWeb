// Package orders предоставляет модель для работы с заказами.
// Данные о заказах хранятся в таблице orders.

package orders

import (
	"time"

	"github.com/Euhcslel/SagaWeb/internal/domain/dealers"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
)

const TableNameOrder = "orders"

type Order struct {
	ID              int64             `gorm:"column:id;primaryKey;autoIncrement:true"`
	DealerID        *int64            `gorm:"column:dealer_id;null;default:NULL"`
	ManagerID       int64             `gorm:"column:manager_id;not null"`
	CreatedAt       time.Time         `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	Status          enums.OrderStatus `gorm:"column:status;type:order_status;not null;"`
	ManufactureDate *time.Time        `gorm:"column:manufacture_date;null;default:NULL"`
	FinalizedAt     *time.Time        `gorm:"column:finalized_at;null;default:NULL"`

	Dealer  *dealers.Dealer `gorm:"foreignKey:DealerID;references:UserID"`
	Manager users.User      `gorm:"foreignKey:ManagerID;references:ID"`
}

func (*Order) TableName() string {
	return TableNameOrder
}
