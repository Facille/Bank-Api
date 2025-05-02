package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/therealadik/bank-api/internal/service"
)

// Ключ для хранения ID пользователя в контексте
type contextKey string

const UserIDKey contextKey = "userID"

// JWTMiddleware middleware для проверки JWT-токена
type JWTMiddleware struct {
	authService service.AuthService
	logger      *logrus.Logger
}

// NewJWTMiddleware создает новый JWT middleware
func NewJWTMiddleware(authService service.AuthService, logger *logrus.Logger) *JWTMiddleware {
	return &JWTMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// Middleware проверяет JWT-токен и добавляет ID пользователя в контекст
func (m *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлечение токена из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		// Формат токена: Bearer <token>
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}

		// Извлечение токена
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		// Проверка и разбор токена
		userID, err := m.authService.ParseToken(tokenString)
		if err != nil {
			m.logger.WithError(err).Warn("Ошибка проверки токена")
			http.Error(w, "Неверный или просроченный токен", http.StatusUnauthorized)
			return
		}

		// Добавление ID пользователя в контекст запроса
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID извлекает ID пользователя из контекста
func GetUserID(ctx context.Context) (int64, error) {
	userID := ctx.Value(UserIDKey).(int64)
	return userID, nil
}
