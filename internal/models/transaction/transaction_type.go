package transaction

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
	Transfer   TransactionType = "transfer"
)
