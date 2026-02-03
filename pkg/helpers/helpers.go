package helpers

import (
	"net/http"
	"project/pkg/database"
	"project/pkg/models"

	"github.com/samborkent/uuidv7"
)

func GetUserBySessionToken(token string) (*models.User, error) {
	var session models.Session
	err := database.DB.
		Preload("User").Preload("User.Role").
		Where("token = ?", token).
		First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session.User, nil
}

func SetSession(w http.ResponseWriter, userId int64) {
	token := uuidv7.New().String()

	session := models.Session{
		UserID:    userId,
		Token: 	token,
	}
	err := database.DB.Create(&session).Error
	if err != nil {
		WriteError(w, err, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}
