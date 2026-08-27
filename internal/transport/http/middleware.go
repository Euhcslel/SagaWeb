package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/types"
	"github.com/Euhcslel/SagaWeb/internal/utils"
	"golang.org/x/time/rate"
)

// RequireAuth проверяет наличие аутентифицированного пользователя в контексте запроса.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.GetUserBySessionToken(r)
		if err != nil {
			helpers.WriteError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// WithOptionalAuth добавляет необязательную аутентификацию к обработчику
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

// CacheControl устанавливает заголовок Cache-Control для всех ответов.
func CacheControl(value string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", value)
			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeaders добавляет базовые security-заголовки ко всем HTTP-ответам.
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

// ipLimiter хранит лимитер и время последнего запроса от IP-адреса.
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mutex    sync.Mutex
	limiters = map[string]*ipLimiter{}
)

// init запускает горутину для очистки лимитеров, которые не использовались более 30 минут.
func init() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			mutex.Lock()
			for ip, l := range limiters {
				if time.Since(l.lastSeen) > 30*time.Minute {
					delete(limiters, ip)
				}
			}
			mutex.Unlock()
		}
	}()
}

// getLimiter вовзращает лимитер по IP-адресу.
func getLimiter(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()
	l, ok := limiters[ip]
	if !ok {
		l = &ipLimiter{limiter: rate.NewLimiter(rate.Every(12*time.Second), 5)}
		limiters[ip] = l
	}
	l.lastSeen = time.Now()
	return l.limiter
}

// RateLimit ограничивает количество запросов от одного IP-адреса.
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
