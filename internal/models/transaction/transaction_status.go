package transaction

type TransactionStatus string

const (
	PENDING   TransactionStatus = "PENDING"
	COMPLETED TransactionStatus = "COMPLETED"
	FAILED    TransactionStatus = "FAILED"
)
