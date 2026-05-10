package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_registration_requests"
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

func CreateDealerRegistrationRequest(request dealer_registration_requests.DealerRegistrationRequest) error {
	return database.DB.Create(&request).Error
}

func DeleteSession(token string) error {
	return database.DB.Where("token = ?", token).Delete(sessions.Session{}).Error
}

func GetAllDealerRegistrationRequests() ([]dealer_registration_requests.DealerRegistrationRequest, error) {
	var requests []dealer_registration_requests.DealerRegistrationRequest
	if err := database.DB.Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}
