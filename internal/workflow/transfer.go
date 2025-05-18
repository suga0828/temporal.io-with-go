package workflow

import (
	"fmt"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"temporal.io-with-go/app/internal/config"
	"temporal.io-with-go/app/internal/domain"
)

// MoneyTransfer orchestrates a money transfer between accounts
func MoneyTransfer(ctx workflow.Context, input domain.TransferDetails) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting money transfer workflow", 
		"sourceAccount", input.SourceAccount,
		"targetAccount", input.TargetAccount, 
		"amount", input.Amount,
		"referenceID", input.ReferenceID)
	
	// Get workflow configuration
	cfg := config.FromEnv().Workflow
	
	// Configure retry policy
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:        cfg.RetryInitialInterval,
		BackoffCoefficient:     cfg.RetryBackoffCoeff,
		MaximumInterval:        cfg.RetryMaxInterval,
		MaximumAttempts:        int32(cfg.RetryMaxAttempts),
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}
	
	// Configure activity options
	options := workflow.ActivityOptions{
		StartToCloseTimeout: cfg.ActivityTimeout,
		RetryPolicy:         retryPolicy,
	}
	
	// Apply the options to context
	ctx = workflow.WithActivityOptions(ctx, options)
	
	// Execute the transfer
	result, err := executeTransfer(ctx, input)
	if err != nil {
		logger.Error("Money transfer failed", "error", err.Error())
		return "", err
	}
	
	logger.Info("Money transfer completed", "result", result)
	return result, nil
}

// executeTransfer handles the actual transfer logic
func executeTransfer(ctx workflow.Context, input domain.TransferDetails) (string, error) {
	logger := workflow.GetLogger(ctx)
	
	logger.Info("Starting money transfer workflow", 
		"sourceAccount", input.SourceAccount, 
		"targetAccount", input.TargetAccount, 
		"amount", input.Amount)
	
	// Step 1: Withdraw money from source account
	// Log step removed - visible in Temporal UI
	logger.Info("Withdrawing money", "account", input.SourceAccount, "amount", input.Amount)
	var withdrawOutput string
	withdrawErr := workflow.ExecuteActivity(ctx, Withdraw, input).Get(ctx, &withdrawOutput)
	
	if withdrawErr != nil {
		logger.Error("Withdrawal failed", "error", withdrawErr.Error())
		return "", fmt.Errorf("failed to withdraw money from %s: %w", 
			input.SourceAccount, withdrawErr)
	}
	
	logger.Info("Withdrawal successful", "transactionID", withdrawOutput)
	
	// Step 2: Deposit money to target account
	// Log step removed - visible in Temporal UI
	logger.Info("Depositing money", "account", input.TargetAccount, "amount", input.Amount)
	var depositOutput string
	depositErr := workflow.ExecuteActivity(ctx, Deposit, input).Get(ctx, &depositOutput)
	
	if depositErr != nil {
		// The deposit failed; attempt to refund the source account
		logger.Error("Deposit failed, initiating refund", "error", depositErr.Error())
		
		var refundOutput string
		refundErr := workflow.ExecuteActivity(ctx, Refund, input).Get(ctx, &refundOutput)
		
		if refundErr != nil {
			// Critical error: both deposit and refund failed
			logger.Error("Refund failed", "error", refundErr.Error())
			return "", fmt.Errorf("deposit to %s failed: %v. refund to %s also failed: %w",
				input.TargetAccount, depositErr, input.SourceAccount, refundErr)
		}
		
		// Deposit failed but refund succeeded
		logger.Info("Refund successful", "transactionID", refundOutput)
		return "", fmt.Errorf("deposit to %s failed: %w. money returned to %s",
			input.TargetAccount, depositErr, input.SourceAccount)
	}
	
	logger.Info("Deposit successful", "transactionID", depositOutput)
	
	// Transfer completed successfully
	result := fmt.Sprintf("Transfer complete (withdrawal: %s, deposit: %s)", 
		withdrawOutput, depositOutput)
	return result, nil
}
