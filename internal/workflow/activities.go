package workflow

import (
	"context"
	"fmt"

	"temporal.io-with-go/app/internal/banking"
	"temporal.io-with-go/app/internal/banking/mock"
	"temporal.io-with-go/app/internal/config"
	"temporal.io-with-go/app/internal/domain"
	"temporal.io-with-go/app/pkg/logger"
)

// getBankingClient returns a configured banking client
func getBankingClient() banking.Client {
	cfg := config.FromEnv()
	return mock.NewClient(cfg.Banking.Hostname)
}

// Withdraw handles the withdrawal activity
func Withdraw(ctx context.Context, data domain.TransferDetails) (string, error) {
	log := logger.New().
		WithField("activity", "Withdraw").
		WithField("sourceAccount", data.SourceAccount).
		WithField("amount", data.Amount).
		WithField("referenceID", data.ReferenceID)
	
	log.Info("Withdrawal started")
	
	// Create a unique reference ID for this withdrawal
	referenceID := fmt.Sprintf("%s-withdrawal", data.ReferenceID)
	
	// Get banking client and execute withdrawal
	client := getBankingClient()
	confirmation, err := client.Withdraw(ctx, data.SourceAccount, data.Amount, referenceID)
	
	if err != nil {
		log.Error(err, "Withdrawal failed", "account", data.SourceAccount)
		return "", err
	}
	
	log.Info("Withdrawal completed")
	return confirmation, nil
}

// Deposit handles the deposit activity
func Deposit(ctx context.Context, data domain.TransferDetails) (string, error) {
	log := logger.New().
		WithField("activity", "Deposit").
		WithField("targetAccount", data.TargetAccount).
		WithField("amount", data.Amount).
		WithField("referenceID", data.ReferenceID)
	
	log.Info("Deposit started")
	
	// Create a unique reference ID for this deposit
	referenceID := fmt.Sprintf("%s-deposit", data.ReferenceID)
	
	// Get banking client and execute deposit
	client := getBankingClient()
	
	// Uncomment the following line to simulate a deposit failure
	// For testing the compensation path, you can use this instead of the normal deposit
	// confirmation, err := client.(*mock.BankingClient).DepositThatFails(ctx, data.TargetAccount, data.Amount, referenceID)
	
	confirmation, err := client.Deposit(ctx, data.TargetAccount, data.Amount, referenceID)
	
	if err != nil {
		log.Error(err, "Deposit failed", "account", data.TargetAccount)
		return "", err
	}
	
	log.Info("Deposit completed")
	return confirmation, nil
}

// Refund handles the refund activity (compensation for failed deposit)
func Refund(ctx context.Context, data domain.TransferDetails) (string, error) {
	log := logger.New().
		WithField("activity", "Refund").
		WithField("sourceAccount", data.SourceAccount).
		WithField("amount", data.Amount).
		WithField("referenceID", data.ReferenceID)
	
	log.Info("Refund started")
	
	// Create a unique reference ID for this refund
	referenceID := fmt.Sprintf("%s-refund", data.ReferenceID)
	
	// Get banking client and execute refund (which is a deposit back to the source account)
	client := getBankingClient()
	confirmation, err := client.Deposit(ctx, data.SourceAccount, data.Amount, referenceID)
	
	if err != nil {
		log.Error(err, "Refund failed", "account", data.SourceAccount, "amount", data.Amount)
		return "", err
	}
	
	log.Info("Refund completed")
	return confirmation, nil
}
