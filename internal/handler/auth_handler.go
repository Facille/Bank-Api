package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Facille/Bank-Api/internal/dto"
	"github.com/Facille/Bank-Api/internal/service"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authService service.AuthService
	logger      *logrus.Logger
}

func NewAuthHandler(authService service.AuthService, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Warn("Ошибка декодирования запроса регистрации")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	userID, err := h.authService.Register(r.Context(), req)
	if err != nil {
		h.logger.WithError(err).Warn("Ошибка при регистрации пользователя")

		if errors.Is(err, service.ErrUserExists) {
			http.Error(w, "Пользователь с таким email или username уже существует", http.StatusConflict)
			return
		}

		http.Error(w, "Ошибка при регистрации пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "Пользователь успешно зарегистрирован",
		"user_id": userID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithError(err).Error("Ошибка при формировании ответа")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Warn("Ошибка декодирования запроса авторизации")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), req)
	if err != nil {
		h.logger.WithError(err).Warn("Ошибка при авторизации пользователя")

		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := dto.AuthResponse{Token: token}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithError(err).Error("Ошибка при формировании ответа авторизации")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
}
