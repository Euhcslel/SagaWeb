package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/companies"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_manager_assignments"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_registration_requests"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealers"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	"github.com/Euhcslel/SagaWeb/internal/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUserDealers(db *gorm.DB, userID int64) ([]dealer_manager_assignments.DealerManagerAssignment, error) {
	var dealers []dealer_manager_assignments.DealerManagerAssignment
	if err := db.Model(dealer_manager_assignments.DealerManagerAssignment{}).Preload("Dealer").Where("manager_id = ?", userID).Find(&dealers).Error; err != nil {
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

func GetDealerByID(db *gorm.DB, dealerID int64) (dealers.Dealer, error) {
	var dealer dealers.Dealer
	if err := db.Where("user_id = ?", dealerID).First(&dealer).Error; err != nil {
		return dealers.Dealer{}, err
	}

	return dealer, nil
}

func UpdateCompanyInfo(db *gorm.DB, companyID int64, userInfo types.UpdatedUserInfo) error {
	return db.Model(&companies.Company{}).
		Where("id = ?", companyID).
		Updates(map[string]any{
			"name": userInfo.Company,
		}).Error
}

func GetRegistrationRequestByID(db *gorm.DB, requestID int64) (dealer_registration_requests.DealerRegistrationRequest, error) {
	var request dealer_registration_requests.DealerRegistrationRequest
	if err := db.Where("id = ?", requestID).First(&request).Error; err != nil {
		return dealer_registration_requests.DealerRegistrationRequest{}, err
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

func AttachDealerToManager(db *gorm.DB, dealerID int64, managerID int64) error {
	return db.Create(&dealer_manager_assignments.DealerManagerAssignment{DealerID: dealerID, ManagerID: managerID}).Error
}

func DeleteRegistrationRequest(db *gorm.DB, requestID int64) error {
	return db.Where("id = ?", requestID).Delete(dealer_registration_requests.DealerRegistrationRequest{}).Error
}
