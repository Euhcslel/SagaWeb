package helpers

import (
	"net/http"
	"project/internal/database"
	"project/internal/domain/sessions"
	"project/internal/domain/users"

	"github.com/samborkent/uuidv7"
)

func GetUserBySessionToken(r *http.Request) (*users.User, error) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	token := sessionToken.Value

	var session sessions.Session
	err = database.DB.
		Preload("User").Preload("User.Role").
		Where("token = ?", token).
		First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session.User, nil
}

func SetSession(w http.ResponseWriter, userId int64) error {
	token := uuidv7.New().String()

	session := sessions.Session{
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
