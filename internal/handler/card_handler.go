package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/therealadik/bank-api/internal/dto"
	"github.com/therealadik/bank-api/internal/middleware"
	"github.com/therealadik/bank-api/internal/service"
)

type CardHandler struct {
	cardService *service.CardService
	logger      *logrus.Logger
}

func NewCardHandler(cardService *service.CardService, logger *logrus.Logger) *CardHandler {
	return &CardHandler{
		cardService: cardService,
		logger:      logger,
	}
}

// CreateCard обработчик для выпуска новой карты
func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		h.logger.Errorf("Ошибка получения userID из контекста: %v", err)
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	// Декодируем запрос
	var req dto.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("Ошибка декодирования запроса: %v", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверяем наличие PGP ключа
	if req.PGPKey == "" {
		h.logger.Warn("Отсутствует PGP ключ")
		http.Error(w, "PGP ключ обязателен", http.StatusBadRequest)
		return
	}

	// Создаем карту
	card, cardDetails, err := h.cardService.CreateCard(r.Context(), userID, req.PGPKey)
	if err != nil {
		h.logger.Errorf("Ошибка создания карты: %v", err)
		http.Error(w, "Не удалось создать карту", http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	resp := dto.CreateCardResponse{
		ID:         card.ID,
		UserID:     card.UserID,
		CreatedAt:  card.CreatedAt.Format("2006-01-02T15:04:05Z"),
		CardNumber: cardDetails["number"],
		Expire:     cardDetails["expire"],
		CVV:        cardDetails["cvv"],
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorf("Ошибка кодирования ответа: %v", err)
	}
}

// GetCards обработчик для получения списка карт пользователя
func (h *CardHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		h.logger.Errorf("Ошибка получения userID из контекста: %v", err)
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	// Получаем список карт
	cards, err := h.cardService.GetUserCards(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("Ошибка получения списка карт: %v", err)
		http.Error(w, "Не удалось получить список карт", http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	resp := dto.CardListResponse{
		Cards: make([]dto.CardResponse, 0, len(cards)),
	}

	for _, card := range cards {
		resp.Cards = append(resp.Cards, dto.CardResponse{
			ID:        card.ID,
			UserID:    card.UserID,
			CreatedAt: card.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorf("Ошибка кодирования ответа: %v", err)
	}
}

// GetCardDetails обработчик для получения детальной информации о карте
func (h *CardHandler) GetCardDetails(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		h.logger.Errorf("Ошибка получения userID из контекста: %v", err)
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	// Получаем ID карты из URL
	vars := mux.Vars(r)
	cardID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		h.logger.Warnf("Неверный формат ID карты: %v", err)
		http.Error(w, "Неверный ID карты", http.StatusBadRequest)
		return
	}

	// Получаем PGP ключ из запроса
	pgpKey := r.URL.Query().Get("pgp_key")
	if pgpKey == "" {
		h.logger.Warn("Отсутствует PGP ключ в запросе")
		http.Error(w, "PGP ключ обязателен", http.StatusBadRequest)
		return
	}

	// Получаем детали карты
	cardDetails, err := h.cardService.GetCardDetails(r.Context(), cardID, userID, pgpKey)
	if err != nil {
		h.logger.Errorf("Ошибка получения данных карты: %v", err)
		http.Error(w, "Не удалось получить данные карты", http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	resp := dto.CardDetailsResponse{
		ID:         cardID,
		CardNumber: cardDetails["number"],
		Expire:     cardDetails["expire"],
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorf("Ошибка кодирования ответа: %v", err)
	}
}

// ProcessPayment обработчик для обработки платежа по карте
func (h *CardHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// Для платежа не требуется быть владельцем карты, только корректные данные карты

	// Декодируем запрос
	var req dto.CardPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("Ошибка декодирования запроса: %v", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if req.CardID == 0 || req.CVV == "" || req.Amount == "" || req.PGPKey == "" {
		h.logger.Warn("Отсутствуют обязательные поля")
		http.Error(w, "Все поля обязательны", http.StatusBadRequest)
		return
	}

	// Проверяем данные карты для платежа
	isValid, err := h.cardService.VerifyCardPayment(r.Context(), req.CardID, req.CVV, req.PGPKey)
	if err != nil {
		h.logger.Errorf("Ошибка проверки данных карты: %v", err)
		http.Error(w, "Ошибка проверки данных карты", http.StatusBadRequest)
		return
	}

	if !isValid {
		h.logger.Warn("Неверные данные карты")
		http.Error(w, "Неверные данные карты", http.StatusBadRequest)
		return
	}

	paymentID := strconv.FormatInt(time.Now().UnixNano(), 10)
	resp := dto.CardPaymentResponse{
		Success:     true,
		PaymentID:   paymentID,
		Description: "Платеж успешно обработан",
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorf("Ошибка кодирования ответа: %v", err)
	}
}
