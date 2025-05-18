package config

import (
	"os"
	"strconv"
	"time"
)

// BankingConfig holds banking service configuration
type BankingConfig struct {
	Hostname string
	Timeout  time.Duration
}

// WorkflowConfig holds Temporal workflow configuration
type WorkflowConfig struct {
	TaskQueue            string
	RetryInitialInterval time.Duration
	RetryBackoffCoeff    float64
	RetryMaxInterval     time.Duration
	RetryMaxAttempts     int
	ActivityTimeout      time.Duration
}

// Config holds the application configuration
type Config struct {
	Banking  BankingConfig
	Workflow WorkflowConfig
}

// FromEnv loads configuration from environment variables
func FromEnv() *Config {
	return &Config{
		Banking: BankingConfig{
			Hostname: getEnvString("BANKING_HOSTNAME", "bank-api.example.com"),
			Timeout:  getEnvDuration("BANKING_TIMEOUT", 10*time.Second),
		},
		Workflow: WorkflowConfig{
			TaskQueue:            getEnvString("WORKFLOW_TASK_QUEUE", "TRANSFER_MONEY_TASK_QUEUE"),
			RetryInitialInterval: getEnvDuration("WORKFLOW_RETRY_INITIAL_INTERVAL", time.Second),
			RetryBackoffCoeff:    getEnvFloat("WORKFLOW_RETRY_BACKOFF_COEFF", 2.0),
			RetryMaxInterval:     getEnvDuration("WORKFLOW_RETRY_MAX_INTERVAL", 100*time.Second),
			RetryMaxAttempts:     getEnvInt("WORKFLOW_RETRY_MAX_ATTEMPTS", 500),
			ActivityTimeout:      getEnvDuration("WORKFLOW_ACTIVITY_TIMEOUT", time.Minute),
		},
	}
}

// Helper functions to get environment variables with defaults

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
