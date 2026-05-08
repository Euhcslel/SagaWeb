package sales

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	"time"
)

const TableNameSale = "sales"

type Sale struct {
	ID        int64             `gorm:"column:id;primaryKey;autoIncrement:true"`
	ClientID  *int64            `gorm:"column:client_id;null;default:NULL"`
	ManagerID int64             `gorm:"column:manager_id;not null"`
	CreatedAt time.Time         `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	Status    enums.OrderStatus `gorm:"column:status;type:order_status;not null;"`

	Client  *users.User `gorm:"foreignKey:ClientID;references:ID"`
	Manager users.User  `gorm:"foreignKey:ManagerID;references:ID"`
}

func (*Sale) TableName() string {
	return TableNameSale
}
