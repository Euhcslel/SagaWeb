package models

const TableNameDealerRegRequest = "dealers_reg_requests"

type DealerRegRequest struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Company     string `gorm:"column:company;not null"`
	Fullname    string `gorm:"column:fullname;not null"`
	Email       string `gorm:"column:email;not null"`
	PhoneNumber string `gorm:"column:phone_number;not null"`
	Password    []byte `gorm:"column:password_hash;not null"`
	StatusID    int32  `gorm:"column:status_id;not null"`

	Status Status `gorm:"foreignKey:StatusID;references:ID"`
}

func (*DealerRegRequest) TableName() string {
	return TableNameDealerRegRequest
}
