package service

import (
	"project/internal/database"
	"project/internal/domain/enums"
	"project/internal/domain/managers_and_dealers"
	"project/internal/domain/users"
	errs "project/internal/errors"
	"project/internal/repository"
	"project/internal/types"
)

func GetUserDealers(user *users.User) ([]managers_and_dealers.ManagerAndDealer, error) {
	role := user.Role

	if role == enums.ManagerRole {
		dealers, err := repository.GetUserDealers(database.DB, user.ID)
		if err != nil {
			return nil, err
		}
		return dealers, nil
	} else {
		return nil, errs.ErrForbidden
	}
}

type UserInfo struct {
	ID       int64
	Fullname string
	Email    string
	Phone    string
	Role     enums.Role
	IsDealer bool

	Dealer *DealerInfo
}

type DealerInfo struct {
	CompanyName string
	Address     string
}

func GetUserInfo(user *users.User) (*UserInfo, error) {
	userInfo := &UserInfo{
		ID:       user.ID,
		Fullname: user.Fullname,
		Email:    user.Email,
		Phone:    user.PhoneNumber,
		Role:     user.Role,
		IsDealer: user.Role == enums.DealerRole,
	}

	if userInfo.IsDealer {
		dealer, err := repository.GetDealerInfo(database.DB, user.ID)
		if err != nil {
			return nil, err
		}

		userInfo.Dealer = &DealerInfo{
			CompanyName: dealer.Company.Name,
			Address:     dealer.Address,
		}
	}

	return userInfo, nil
}

func UpdateUserInfo(user *users.User, userInfo types.UpdatedUserInfo) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := repository.UpdateUserInfo(tx, user.ID, userInfo); err != nil {
		tx.Rollback()
		return err
	}

	if user.Role == enums.DealerRole {
		dealer, err := repository.GetDealerById(tx, user.ID)
		if err != nil {
			tx.Rollback()

			return err
		}

		if err := repository.UpdateDealerInfo(tx, dealer.UserID, userInfo); err != nil {
			tx.Rollback()
			return err
		}

		if err := repository.UpdateCompanyInfo(tx, dealer.CompanyID, userInfo); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
