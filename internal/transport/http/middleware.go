package http

import (
	"context"
	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/types"
	"github.com/Euhcslel/SagaWeb/internal/utils"
	"net/http"
)

// Middleware для проверки аутентификации
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.GetUserBySessionToken(r)
		if err != nil {
			helpers.WriteError(w, err, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// Middleware для необязательной аутентификации
func WithOptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := utils.GetUserBySessionToken(r)
		if user != nil {
			ctx := context.WithValue(r.Context(), types.UserContextKey, user)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// Middleware, который добавляет базовые security-заголовки ко всем HTTP-ответам
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
