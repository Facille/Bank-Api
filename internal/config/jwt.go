package config

import (
	"os"
	"time"
)

// JWTConfig содержит настройки для JWT-токенов
type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

// LoadJWT загружает конфигурацию JWT из переменных окружения
func LoadJWT() JWTConfig {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-bank-api-jwt-secret-key" // В продакшене не использовать дефолтный ключ!
	}

	// TTL токена 24 часа как указано в ТЗ
	return JWTConfig{
		Secret:    secret,
		ExpiresIn: 24 * time.Hour,
	}
}
