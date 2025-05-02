package transaction

type TransactionType string

const (
	DEPOSIT    TransactionType = "DEPOSIT"
	WITHDRAWAL TransactionType = "WITHDRAWAL"
	TRANSFER   TransactionType = "TRANSFER"
)
