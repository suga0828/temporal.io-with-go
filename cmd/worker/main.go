package main

import (
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"temporal.io-with-go/app/internal/config"
	"temporal.io-with-go/app/internal/workflow"
	"temporal.io-with-go/app/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.New().WithField("component", "worker")
	log.Info("Starting money transfer worker")

	// Load configuration
	cfg := config.FromEnv()
	
	// Get Temporal server address from environment variable or use default
	temporalAddress := os.Getenv("TEMPORAL_ADDRESS")
	if temporalAddress == "" {
		temporalAddress = client.DefaultHostPort
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

	// Create worker
	w := worker.New(c, cfg.Workflow.TaskQueue, worker.Options{})

	// Register workflow and activities
	w.RegisterWorkflow(workflow.MoneyTransfer)
	w.RegisterActivity(workflow.Withdraw)
	w.RegisterActivity(workflow.Deposit)
	w.RegisterActivity(workflow.Refund)

	// Start listening to the task queue
	log.Info("Worker started", "taskQueue", cfg.Workflow.TaskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Error(err, "Worker execution failed")
		os.Exit(1)
	}
}
