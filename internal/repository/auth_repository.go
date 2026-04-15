package repository

import (
	"project/internal/database"
	"project/internal/domain/dealers_reg_requests"
	"project/internal/domain/sessions"
	"project/internal/domain/users"
)

func GetUserByUsername(username string) (users.User, error) {
	var user users.User
	if err := database.DB.
		Where("username = ?", username).
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
