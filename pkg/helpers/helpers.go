package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"project/pkg/database"
	"project/pkg/models"
)

func GetUserBySessionToken(token string) models.User {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])
	var session models.Session
	err := database.DB.
		Preload("User").Preload("User.Role").
		Where("token_hash = ?", tokenHash).
		First(&session).Error
	if err != nil {
		return models.User{}
	}
	return session.User
}

func SetSession(w http.ResponseWriter, userId int64) {
	// Генерация токена
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		WriteErrorDebug(w, err, http.StatusInternalServerError)
		return
	}
	// Хэширование для куки
	token := base64.RawURLEncoding.EncodeToString(b)

	// Хэширование для БД
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	session := models.Session{
		UserID:    userId,
		TokenHash: tokenHash,
	}
	err := database.DB.Create(&session).Error
	if err != nil {
		WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}
