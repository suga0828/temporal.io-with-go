package banking

import "context"

// Client defines the interface for banking operations
type Client interface {
	// Withdraw takes money from an account
	// Returns a transaction ID on success, or an error
	Withdraw(ctx context.Context, accountNumber string, amount int, referenceID string) (string, error)

	// Deposit adds money to an account
	// Returns a transaction ID on success, or an error
	Deposit(ctx context.Context, accountNumber string, amount int, referenceID string) (string, error)
}
