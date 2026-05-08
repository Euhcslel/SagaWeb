package utils

import (
	"context"
	"net/http"
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/sessions"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	"github.com/Euhcslel/SagaWeb/internal/types"

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
		Preload("User").
		Where("token = ?", token).
		First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session.User, nil
}

func SetSession(w http.ResponseWriter, userID int64) error {
	token := uuidv7.New().String()

	session := sessions.Session{
		UserID: userID,
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

// Функция для получения пользователя из контекста
func UserFromContext(ctx context.Context) *users.User {
	user, ok := ctx.Value(types.UserContextKey).(*users.User)
	if !ok {
		return nil
	}
	return user
}

func HasDealerAccess(role enums.Role) bool {
	return role == enums.DealerRole || role == enums.AdminRole || role == enums.ManagerRole || role == enums.LogisticianRole
}
