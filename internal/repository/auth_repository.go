package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealers_reg_requests"
	"github.com/Euhcslel/SagaWeb/internal/domain/sessions"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
)

func GetUserByEmail(email string) (users.User, error) {
	var user users.User
	if err := database.DB.
		Where("email = ?", email).
		First(&user).
		Error; err != nil {
		return users.User{}, err
	}
	return user, nil
}

func CreateDealerRegRequest(request dealers_reg_requests.DealerRegRequest) error {
	return database.DB.Create(&request).Error
}

func DeleteSession(token string) error {
	return database.DB.Where("token = ?", token).Delete(sessions.Session{}).Error
}

func GetAllDealerRegRequests() ([]dealers_reg_requests.DealerRegRequest, error) {
	var requests []dealers_reg_requests.DealerRegRequest
	if err := database.DB.Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}
