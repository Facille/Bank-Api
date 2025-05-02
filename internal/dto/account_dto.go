package dto

import (
	"github.com/shopspring/decimal"
	"github.com/therealadik/bank-api/internal/models/account"
	"github.com/therealadik/bank-api/internal/models/transaction"
)

// CreateAccountRequest запрос на создание счета
type CreateAccountRequest struct {
	Currency account.Currency `json:"currency"`
}

// UpdateBalanceRequest запрос на пополнение/списание
type UpdateBalanceRequest struct {
	Amount decimal.Decimal `json:"amount"`
}

// TransferRequest запрос на перевод между счетами
type TransferRequest struct {
	FromAccountID int64           `json:"from_account_id"`
	ToAccountID   int64           `json:"to_account_id"`
	Amount        decimal.Decimal `json:"amount"`
}

// AccountResponse ответ со счетом
type AccountResponse struct {
	ID        int64            `json:"id"`
	UserID    int64            `json:"user_id"`
	Balance   decimal.Decimal  `json:"balance"`
	Currency  account.Currency `json:"currency"`
	CreatedAt string           `json:"created_at"`
}

// TransactionResponse ответ с транзакцией
type TransactionResponse struct {
	ID        int64              `json:"id"`
	AccountID int64              `json:"account_id"`
	Amount    decimal.Decimal    `json:"amount"`
	Type      transaction.Type   `json:"type"`
	Status    transaction.Status `json:"status"`
	CreatedAt string             `json:"created_at"`
}

// AccountsListResponse список счетов
type AccountsListResponse struct {
	Accounts []AccountResponse `json:"accounts"`
}

// TransactionListResponse список транзакций
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
}
