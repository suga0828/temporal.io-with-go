package banking

import "fmt"

// InsufficientFundsError is raised when the account doesn't have enough money
type InsufficientFundsError struct {
	AccountID string
	Available int
	Requested int
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds in account %s: available %d, requested %d",
		e.AccountID, e.Available, e.Requested)
}

// InvalidAccountError is raised when the account number is invalid
type InvalidAccountError struct {
	AccountID string
}

func (e *InvalidAccountError) Error() string {
	return fmt.Sprintf("account number %s is invalid", e.AccountID)
}

// TransactionError is raised for general banking transaction errors
type TransactionError struct {
	Operation   string
	AccountID   string
	ReferenceID string
	Message     string
}

func (e *TransactionError) Error() string {
	return fmt.Sprintf("%s failed for account %s (ref: %s): %s",
		e.Operation, e.AccountID, e.ReferenceID, e.Message)
}
