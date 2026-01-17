package models

const TableNameDealerRegRequest = "dealers_reg_requests"

type DealerRegRequest struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Company     string `gorm:"column:company;not null" json:"company"`
	Fullname    string `gorm:"column:fullname;not null" json:"fullname"`
	Email       string `gorm:"column:email;not null" json:"email"`
	PhoneNumber string `gorm:"column:phone_number;not null" json:"phone_number"`
	StatusID    int64  `gorm:"column:status_id;not null" json:"status_id"`

	Status Status `gorm:"foreignKey:StatusID;references:ID"`
}

func (*DealerRegRequest) TableName() string {
	return TableNameDealerRegRequest
}
