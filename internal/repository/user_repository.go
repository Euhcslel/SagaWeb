package repository

import (
	"project/internal/domain/dealers"
	"project/internal/domain/managers_and_dealers"

	"gorm.io/gorm"
)

func GetUserDealers(db *gorm.DB, userID int64) ([]managers_and_dealers.ManagerAndDealer, error) {
	var dealers []managers_and_dealers.ManagerAndDealer
	if err := db.Model(managers_and_dealers.ManagerAndDealer{}).Preload("Dealer").Where("manager_id = ?", userID).Find(&dealers).Error; err != nil {
		return nil, err
	}

	return dealers, nil
}

func GetDealerInfo(db *gorm.DB, userID int64) (*dealers.Dealer, error) {
	var dealer dealers.Dealer
	if err := db.Preload("Company").Where("user_id = ?", userID).First(&dealer).Error; err != nil {
		return nil, err
	}
	return &dealer, nil
}
