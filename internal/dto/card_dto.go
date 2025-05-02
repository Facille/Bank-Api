package dto

// CreateCardRequest запрос на создание новой карты
type CreateCardRequest struct {
	PGPKey string `json:"pgp_key"`
}

// CreateCardResponse ответ с данными созданной карты
type CreateCardResponse struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	CreatedAt  string `json:"created_at"`
	CardNumber string `json:"card_number"`
	Expire     string `json:"expire"`
	CVV        string `json:"cvv"`
}

// CardResponse ответ с базовыми данными карты (без секретных данных)
type CardResponse struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

// CardDetailsResponse ответ с деталями карты
type CardDetailsResponse struct {
	ID         int64  `json:"id"`
	CardNumber string `json:"card_number"` // Маскированный номер
	Expire     string `json:"expire"`
}

// CardListResponse список карт
type CardListResponse struct {
	Cards []CardResponse `json:"cards"`
}

// CardPaymentRequest запрос на оплату картой
type CardPaymentRequest struct {
	CardID int64  `json:"card_id"`
	Amount string `json:"amount"`
	CVV    string `json:"cvv"`
	PGPKey string `json:"pgp_key"`
}

// CardPaymentResponse ответ на запрос оплаты
type CardPaymentResponse struct {
	Success     bool   `json:"success"`
	PaymentID   string `json:"payment_id,omitempty"`
	Description string `json:"description,omitempty"`
}
