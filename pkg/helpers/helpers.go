package helpers

import (
	"net/http"
	"project/pkg/database"
	"project/pkg/models"
	"time"

	"github.com/samborkent/uuidv7"
)

func GetUserBySessionToken(w http.ResponseWriter, r *http.Request) *models.User {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	token := sessionToken.Value

	var session models.Session
	err = database.DB.
		Preload("User").Preload("User.Role").
		Where("token = ?", token).
		First(&session).Error
	if err != nil {
		WriteError(w, err, http.StatusInternalServerError)
		return nil
	} else if session.ExpiresAt.Before(time.Now()) {
		if err := database.DB.Delete(models.Session{}).Where("token = ?", token).Error; err != nil {
			WriteError(w, err, http.StatusInternalServerError)
		}
		return nil
	}

	return &session.User
}

func SetSession(w http.ResponseWriter, userId int64) error {
	token := uuidv7.New().String()

	session := models.Session{
		UserID: userId,
		Token:  token,
	}
	if err := database.DB.Create(&session).Error; err != nil {
		return err
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
	return nil
}
