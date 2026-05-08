package dealers_reg_requests

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
)

const TableNameDealerRegRequest = "dealers_reg_requests"

type DealerRegRequest struct {
	ID           int64                  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Company      string                 `gorm:"column:company;not null"`
	Fullname     string                 `gorm:"column:fullname;not null"`
	Email        string                 `gorm:"column:email;not null"`
	PhoneNumber  string                 `gorm:"column:phone_number;not null"`
	Status       enums.RegRequestStatus `gorm:"type:reg_request_status;not null;default:'pending'"`
}

func (*DealerRegRequest) TableName() string {
	return TableNameDealerRegRequest
}
