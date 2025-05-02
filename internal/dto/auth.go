package dto

// RegisterRequest - запрос на регистрацию пользователя
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest - запрос на аутентификацию пользователя
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse - ответ с токеном аутентификации
type AuthResponse struct {
	Token string `json:"token"`
}
