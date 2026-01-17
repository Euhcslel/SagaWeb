package models

const TableNameManagerAndDealer = "managers_and_dealers"

type ManagerAndDealer struct {
	ManagerID int64 `gorm:"column:manager_id;primaryKey" json:"manager_id"`
	DealerID  int64 `gorm:"column:dealer_id;primaryKey" json:"dealer_id"`

	Manager User   `gorm:"foreignKey:ManagerID;references:ID"`
	Dealer  User `gorm:"foreignKey:DealerID;references:ID"`
}

func (*ManagerAndDealer) TableName() string {
	return TableNameManagerAndDealer
}
