package credit

type CreditStatus string

const (
	CreditStatusActive  CreditStatus = "active"
	CreditStatusClosed  CreditStatus = "closed"
	CreditStatusOverdue CreditStatus = "overdue"
)
