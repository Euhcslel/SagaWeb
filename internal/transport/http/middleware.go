package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/types"
	"github.com/Euhcslel/SagaWeb/internal/utils"
	"golang.org/x/time/rate"
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

// Middleware для управления кэшированием HTTP-ответов
func CacheControl(value string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", value)
			next.ServeHTTP(w, r)
		})
	}
}

// Middleware, который добавляет базовые security-заголовки ко всем HTTP-ответам
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-src https://yandex.ru https://api-maps.yandex.ru; frame-ancestors 'none'")
		next.ServeHTTP(w, r)
	})
}

var (
	mutex    sync.Mutex
	limiters = map[string]*rate.Limiter{}
)

// Функция для получения лимитера по IP-адресу
func getLimiter(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()
	limit, ok := limiters[ip]
	if !ok {
		limit = rate.NewLimiter(rate.Every(12*time.Second), 5)
		limiters[ip] = limit
	}
	return limit
}

// Middleware, ограничивающий количество запросов от одного IP-адреса
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if !getLimiter(ip).Allow() {
			helpers.WriteError(w, errors.New("Cлишком много попыток. Подождите."), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
