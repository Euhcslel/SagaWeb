package service

import (
	"project/internal/database"
	"project/internal/domain/managers_and_dealers"
	"project/internal/domain/users"
	errs "project/internal/errors"
	"project/internal/repository"
)

func GetUserDealers(user *users.User) ([]managers_and_dealers.ManagerAndDealer, error) {
	role := user.Role.Name

	if role == "manager" {
		dealers, err := repository.GetUserDealers(database.DB, user.ID)
		if err != nil {
			return nil, err
		}
		return dealers, nil
	} else {
		return nil, errs.ErrForbidden
	}
}
