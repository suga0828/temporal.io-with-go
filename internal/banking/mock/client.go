package mock

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"temporal.io-with-go/app/internal/banking"
)

// Ensure BankingClient implements the banking.Client interface
var _ banking.Client = (*BankingClient)(nil)

// account represents a bank account in the mock system
type account struct {
	AccountNumber string
	Balance       int64
}

// mockBank is our in-memory bank database
type mockBank struct {
	Accounts []account
}

// findAccount looks up an account by account number
func (b mockBank) findAccount(accountNumber string) (account, error) {
	for _, v := range b.Accounts {
		if v.AccountNumber == accountNumber {
			return v, nil
		}
	}
	return account{}, errors.New("account not found")
}

// Initialize our mock bank with some test accounts
var bankDB = &mockBank{
	Accounts: []account{
		{AccountNumber: "85-150", Balance: 2000},
		{AccountNumber: "43-812", Balance: 0},
	},
}

// BankingClient implements the banking.Client interface for testing
type BankingClient struct {
	Hostname string
}

// NewClient creates a new mock banking client
func NewClient(hostname string) *BankingClient {
	return &BankingClient{
		Hostname: hostname,
	}
}

// Withdraw simulates a withdrawal from a bank account
func (client *BankingClient) Withdraw(ctx context.Context, accountNumber string, amount int, referenceID string) (string, error) {
	// Check for context cancellation
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// Find the account
	account, err := bankDB.findAccount(accountNumber)
	if err != nil {
		return "", &banking.InvalidAccountError{AccountID: accountNumber}
	}

	// Check for sufficient funds
	if account.Balance < int64(amount) {
		return "", &banking.InsufficientFundsError{
			AccountID:  accountNumber,
			Available:  int(account.Balance),
			Requested:  amount,
		}
	}

	// Generate a transaction ID and return
	return generateTransactionID("W", 10), nil
}

// Deposit simulates a deposit to a bank account
func (client *BankingClient) Deposit(ctx context.Context, accountNumber string, amount int, referenceID string) (string, error) {
	// Check for context cancellation
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// Find the account
	_, err := bankDB.findAccount(accountNumber)
	if err != nil {
		return "", &banking.InvalidAccountError{AccountID: accountNumber}
	}

	// Generate a transaction ID and return
	return generateTransactionID("D", 10), nil
}

// DepositThatFails simulates a failed deposit operation
func (client *BankingClient) DepositThatFails(ctx context.Context, accountNumber string, amount int, referenceID string) (string, error) {
	return "", &banking.TransactionError{
		Operation:   "deposit",
		AccountID:   accountNumber,
		ReferenceID: referenceID,
		Message:     "simulated deposit failure",
	}
}

// generateTransactionID creates a random transaction ID with the given prefix and length
func generateTransactionID(prefix string, length int) string {
	// Initialize random number generator if needed
	rand.Seed(time.Now().UnixNano())

	// Generate random characters
	randChars := make([]byte, length)
	for i := range randChars {
		allowedChars := "0123456789"
		randChars[i] = allowedChars[rand.Intn(len(allowedChars))]
	}
	return prefix + string(randChars)
}
