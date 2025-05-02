package credit

type CreditStatus string

const (
	ACTIVE  CreditStatus = "ACTIVE"
	CLOSED  CreditStatus = "CLOSED"
	OVERDUE CreditStatus = "OVERDUE"
)
