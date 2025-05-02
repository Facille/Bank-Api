package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/therealadik/bank-api/internal/config"
	"github.com/therealadik/bank-api/internal/dto"
	"github.com/therealadik/bank-api/internal/models"
	"github.com/therealadik/bank-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Различные ошибки, которые могут возникнуть в процессе аутентификации
var (
	ErrInvalidCredentials = errors.New("неверные учетные данные")
	ErrUserExists         = errors.New("пользователь уже существует")
)

// AuthService интерфейс для сервиса аутентификации
type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (int64, error)
	Login(ctx context.Context, req dto.LoginRequest) (string, error)
	ParseToken(tokenString string) (int64, error)
}

// authService реализация сервиса аутентификации
type authService struct {
	userRepo repository.UserRepository
	jwtCfg   config.JWTConfig
}

// NewAuthService создает новый сервис аутентификации
func NewAuthService(userRepo repository.UserRepository, jwtCfg config.JWTConfig) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
	}
}

// Register регистрирует нового пользователя
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
	// Хеширование пароля с использованием bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Login аутентифицирует пользователя и возвращает JWT-токен
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// Генерация JWT-токена
	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateToken генерирует JWT-токен
func (s *authService) generateToken(userID int64) (string, error) {
	// Структура claims для JWT
	claims := jwt.MapClaims{
		"sub": userID,                                    // subject (ID пользователя)
		"exp": time.Now().Add(s.jwtCfg.ExpiresIn).Unix(), // expiration time
		"iat": time.Now().Unix(),                         // issued at
	}

	// Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписание токена
	tokenString, err := token.SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken разбирает и проверяет JWT-токен, возвращает ID пользователя
func (s *authService) ParseToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неожиданный метод подписи токена")
		}
		return []byte(s.jwtCfg.Secret), nil
	})

	if err != nil {
		return 0, err
	}

	// Проверка валидности токена
	if !token.Valid {
		return 0, errors.New("невалидный токен")
	}

	// Извлечение claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("невалидные claims")
	}

	// Извлечение ID пользователя
	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("невалидный ID пользователя")
	}

	return int64(userID), nil
}
