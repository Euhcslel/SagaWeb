package repository

import (
	"project/internal/database"
	"project/internal/domain/companies"
	"project/internal/domain/dealers"
	"project/internal/domain/dealers_reg_requests"
	"project/internal/domain/managers_and_dealers"
	"project/internal/domain/users"
	"project/internal/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
			"fullname":     userInfo.Fullname,
			"email":        userInfo.Email,
			"phone_number": userInfo.Phone,
		}).Error
}

func UpdateDealerInfo(db *gorm.DB, userID int64, userInfo types.UpdatedUserInfo) error {
	return db.Model(&dealers.Dealer{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"address": userInfo.Address,
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
			"name": userInfo.Company,
		}).Error
}

func DeleteDealerRegRequest(requestId int64) error {
	return database.DB.Where("id = ?", requestId).Delete(dealers_reg_requests.DealerRegRequest{}).Error
}

func GetRegRequestById(db *gorm.DB, requestId int64) (dealers_reg_requests.DealerRegRequest, error) {
	var request dealers_reg_requests.DealerRegRequest
	if err := database.DB.Where("id = ?", requestId).First(&request).Error; err != nil {
		return dealers_reg_requests.DealerRegRequest{}, err
	}

	return request, nil
}

func CreateUser(db *gorm.DB, user *users.User) error {
	return db.
		Clauses(clause.Returning{}).
		Create(user).
		Error
}

func CreateDealer(db *gorm.DB, dealer *dealers.Dealer) error {
	return db.
		Clauses(clause.Returning{}).
		Create(dealer).
		Error
}

func GetOrCreateCompanyByName(db *gorm.DB, name string) (*companies.Company, error) {
	var company companies.Company
	err := db.
		Where("name = ?", name).
		FirstOrCreate(&company, companies.Company{
			Name: name,
		}).
		Error

	if err != nil {
		return nil, err
	}

	return &company, nil
}

func AttachDealerToManager(db *gorm.DB, dealerId int64, managerId int64) error {
	return db.Create(&managers_and_dealers.ManagerAndDealer{DealerID: dealerId, ManagerID: managerId}).Error
}

func DeleteRegRequest(db *gorm.DB, requestId int64) error {
	return db.Where("id = ?", requestId).Delete(dealers_reg_requests.DealerRegRequest{}).Error
}
