package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_registration_requests"
	"github.com/Euhcslel/SagaWeb/internal/domain/sessions"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
)

// GetUserByEmail возвращает информацию о пользователе по его email.
func GetUserByEmail(email string) (*users.User, error) {
	var user users.User
	if err := database.DB.
		Where("email = ?", email).
		First(&user).
		Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateDealerRegistrationRequest создает новый запрос на регистрацию дилера в базе данных.
func CreateDealerRegistrationRequest(request dealer_registration_requests.DealerRegistrationRequest) error {
	return database.DB.Create(&request).Error
}

// DeleteSession удаляет информацию о сессии из базы данных по её токену.
func DeleteSession(token string) error {
	return database.DB.Where("token = ?", token).Delete(&sessions.Session{}).Error
}

// GetAllDealerRegistrationRequests возвращает список всех запросов на регистрацию дилеров.
func GetAllDealerRegistrationRequests() ([]dealer_registration_requests.DealerRegistrationRequest, error) {
	var requests []dealer_registration_requests.DealerRegistrationRequest
	if err := database.DB.Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}
