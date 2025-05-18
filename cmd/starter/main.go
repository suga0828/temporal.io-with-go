package main

import (
	"context"
	"os"

	"go.temporal.io/sdk/client"

	"temporal.io-with-go/app/internal/config"
	"temporal.io-with-go/app/internal/domain"
	"temporal.io-with-go/app/internal/workflow"
	"temporal.io-with-go/app/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.New().WithField("component", "starter")
	log.Info("Starting money transfer workflow client")

	// Load configuration
	cfg := config.FromEnv()
	
	// Get Temporal server address from environment variable or use default
	temporalAddress := os.Getenv("TEMPORAL_ADDRESS")
	if temporalAddress == "" {
		temporalAddress = client.DefaultHostPort
	}
	
	// Get workflow ID from environment variable or use default
	workflowID := os.Getenv("WORKFLOW_ID")
	if workflowID == "" {
		workflowID = "money-transfer-workflow"
	}
	
	log.Info("Connecting to Temporal server", "address", temporalAddress)

	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort: temporalAddress,
	})
	if err != nil {
		log.Error(err, "Unable to create Temporal client")
		os.Exit(1)
	}
	defer c.Close()

	// Prepare payment details
	input := domain.TransferDetails{
		SourceAccount: "85-150",
		TargetAccount: "43-812",
		Amount:        250,
		ReferenceID:   "tx-12345",
	}

	// Configure workflow options
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: cfg.Workflow.TaskQueue,
	}

	log.Info("Initiating workflow", 
		"sourceAccount", input.SourceAccount, 
		"targetAccount", input.TargetAccount, 
		"amount", input.Amount, 
		"workflowID", workflowID)

	// Execute the workflow
	we, err := c.ExecuteWorkflow(context.Background(), options, workflow.MoneyTransfer, input)
	if err != nil {
		log.Error(err, "Unable to start workflow")
		os.Exit(1)
	}

	log.Info("Workflow started", "workflowID", we.GetID(), "runID", we.GetRunID())

	// Wait for workflow completion
	var result string
	if err := we.Get(context.Background(), &result); err != nil {
		log.Error(err, "Workflow execution failed")
		os.Exit(1)
	}

	log.Info("Workflow completed successfully", "result", result)
}
