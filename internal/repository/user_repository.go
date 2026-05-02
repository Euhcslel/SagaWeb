package repository

import (
	"project/internal/domain/companies"
	"project/internal/domain/dealers"
	"project/internal/domain/managers_and_dealers"
	"project/internal/domain/users"
	"project/internal/types"

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

func UpdateUserInfo(db *gorm.DB, userID int64, userInfo types.UpdatedUserInfo) error {
	return db.Model(&users.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"fullname": userInfo.Fullname,
			"email":    userInfo.Email,
			"phone_number":    userInfo.Phone,
		}).Error
}

func UpdateDealerInfo(db *gorm.DB, userID int64, userInfo types.UpdatedUserInfo) error {
	return db.Model(&dealers.Dealer{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"address":    userInfo.Address,
		}).Error
}

func GetDealerById(db *gorm.DB, dealerId int64) (dealers.Dealer, error) {
	var dealer dealers.Dealer
	if err := db.Where("user_id = ?", dealerId).First(&dealer).Error; err != nil {
		return dealers.Dealer{}, err
	}

	return dealer, nil
}

func UpdateCompanyInfo(db *gorm.DB, companyId int64, userInfo types.UpdatedUserInfo) error {
	return db.Model(&companies.Company{}).
		Where("id = ?", companyId).
		Updates(map[string]any{
			"name":    userInfo.Company,
		}).Error
}
