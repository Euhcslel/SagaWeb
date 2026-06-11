package repository

import (
	"fmt"

	"github.com/Euhcslel/SagaWeb/internal/domain/companies"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_contract"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_manager_assignments"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_registration_requests"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealers"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	"github.com/Euhcslel/SagaWeb/internal/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUserDealers(db *gorm.DB, userID int64) ([]dealer_manager_assignments.DealerManagerAssignment, error) {
	var result []dealer_manager_assignments.DealerManagerAssignment
	err := db.Model(dealer_manager_assignments.DealerManagerAssignment{}).
		Preload("Dealer").
		Preload("Dealer.User").
		Preload("Dealer.Company").
		Where("manager_id = ?", userID).
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func UpsertDealerContract(db *gorm.DB, contract dealer_contract.DealerContract) error {
	return db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&contract).Error
}

func GetDealerContract(db *gorm.DB, dealerID int64) (*dealer_contract.DealerContract, error) {
	var contract dealer_contract.DealerContract
	res := db.Where("dealer_id = ?", dealerID).Limit(1).Find(&contract)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &contract, nil
}

func GetDealerContractsByDealerIDs(db *gorm.DB, dealerIDs []int64) (map[int64]*dealer_contract.DealerContract, error) {
	result := make(map[int64]*dealer_contract.DealerContract)
	if len(dealerIDs) == 0 {
		return result, nil
	}
	var contracts []dealer_contract.DealerContract
	if err := db.Where("dealer_id IN ?", dealerIDs).Find(&contracts).Error; err != nil {
		return nil, err
	}
	for i := range contracts {
		result[contracts[i].DealerID] = &contracts[i]
	}
	return result, nil
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

func GetDealerByID(db *gorm.DB, dealerID int64) (*dealers.Dealer, error) {
	var dealer dealers.Dealer
	if err := db.Where("user_id = ?", dealerID).First(&dealer).Error; err != nil {
		return nil, err
	}

	return &dealer, nil
}

func UpdateCompanyInfo(db *gorm.DB, companyID int64, userInfo types.UpdatedUserInfo) error {
	return db.Model(&companies.Company{}).
		Where("id = ?", companyID).
		Updates(map[string]any{
			"name": userInfo.Company,
		}).Error
}

func UpdateCompanyDetails(db *gorm.DB, companyID int64, inn, kpp, ogrn string) error {
	return db.Model(&companies.Company{}).
		Where("id = ?", companyID).
		Updates(map[string]any{
			"inn":  inn,
			"kpp":  kpp,
			"ogrn": ogrn,
		}).Error
}

func GetDealerWithUser(db *gorm.DB, dealerID int64) (*dealers.Dealer, error) {
	var dealer dealers.Dealer
	if err := db.Preload("Company").Preload("User").Where("user_id = ?", dealerID).First(&dealer).Error; err != nil {
		return nil, err
	}
	return &dealer, nil
}

// GetNextContractNumber возвращает MAX числового суффикса из contract_number
// вида "contractN" + 1. Если таких записей нет — возвращает 1.
func GetNextContractNumber(db *gorm.DB) (int, error) {
	var numbers []string
	err := db.Model(&dealer_contract.DealerContract{}).
		Where("contract_number ~ '^contract[0-9]+$'").
		Pluck("contract_number", &numbers).Error
	if err != nil {
		return 1, err
	}

	max := 0
	for _, n := range numbers {
		v := 0
		fmt.Sscanf(n, "contract%d", &v)
		if v > max {
			max = v
		}
	}
	return max + 1, nil
}

func GetRegistrationRequestByID(db *gorm.DB, requestID int64) (*dealer_registration_requests.DealerRegistrationRequest, error) {
	var request dealer_registration_requests.DealerRegistrationRequest
	if err := db.Where("id = ?", requestID).First(&request).Error; err != nil {
		return nil, err
	}

	return &request, nil
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
	if err := db.
		Where("name = ?", name).
		FirstOrCreate(&company, companies.Company{
			Name: name,
		}).
		Error; err != nil {
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
