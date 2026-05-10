package dealer_registration_requests

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
)

const TableNameDealerRegistrationRequest = "dealer_registration_requests"

type DealerRegistrationRequest struct {
	ID          int64                           `gorm:"column:id;primaryKey;autoIncrement:true"`
	Company     string                          `gorm:"column:company;not null"`
	Fullname    string                          `gorm:"column:fullname;not null"`
	Email       string                          `gorm:"column:email;not null"`
	PhoneNumber string                          `gorm:"column:phone_number;not null"`
	Status      enums.RegistrationRequestStatus `gorm:"type:registration_request_status;not null;default:'pending'"`
}

func (*DealerRegistrationRequest) TableName() string {
	return TableNameDealerRegistrationRequest
}
