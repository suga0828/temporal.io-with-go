# Temporal Money Transfer Application

This is a production-ready implementation of a money transfer application using Temporal and Go. It demonstrates how to build robust, fault-tolerant workflows with compensation handling for financial transactions.

The application is based on the Temporal Go SDK tutorial: https://learn.temporal.io/getting_started/go/first_program_in_go/

## Project Structure

The application follows a clean architecture approach with clear separation of concerns:

```
temporal-money-transfer/
├── cmd/                         # Command-line applications
│   ├── worker/                  # Worker process
│   └── starter/                 # Workflow starter
├── internal/                    # Private application code
│   ├── config/                  # Configuration management
│   ├── domain/                  # Core domain models
│   ├── banking/                 # Banking service integration
│   │   └── mock/                # Mock implementation
│   └── workflow/                # Temporal workflows
├── pkg/                         # Public libraries
│   └── logger/                  # Logging utilities
└── scripts/                     # Operational scripts
```

## Features

- **Robust Error Handling**: Comprehensive error types and handling
- **Configuration Management**: Environment-based configuration
- **Structured Logging**: Consistent logging across components
- **Clean Architecture**: Separation of concerns with interfaces
- **Compensation Logic**: Automatic refunds on failed deposits
- **Environment Management**: direnv integration for easy setup

## Getting Started

### Prerequisites

- Go 1.16 or later
- [Temporal Server](https://docs.temporal.io/docs/server/quick-install/)
- (Optional) direnv for environment management

### Running with the Convenience Script

The easiest way to run the application is with the provided script:

```bash
# Make the script executable
chmod +x start-temporal.sh

# Run with default settings
./start-temporal.sh

# Or customize settings
TEMPORAL_UI_PORT=8081 WORKFLOW_ID="custom-transfer-123" ./start-temporal.sh
```

The script will:
1. Start the Temporal server
2. Launch the worker
3. Execute the workflow
4. Monitor all processes

### Manual Execution

Alternatively, you can run each component separately:

#### Step 1: Start Temporal Server

```bash
temporal server start-dev --db-filename temporal.db
```

#### Step 2: Run the Worker

```bash
cd cmd/worker && go run main.go
```

#### Step 3: Execute the Workflow

```bash
cd cmd/starter && go run main.go
```

## Configuration

The application uses environment variables for configuration. You can set these in your shell, use the `.env` file, or leverage direnv with the provided `.envrc` file.

Key configuration options:

```
# Temporal server settings
TEMPORAL_UI_PORT=8080
TEMPORAL_PORT=7233
TEMPORAL_DB_FILE="temporal.db"

# Workflow settings
WORKFLOW_ID="money-transfer-workflow"
WORKFLOW_TASK_QUEUE="TRANSFER_MONEY_TASK_QUEUE"
WORKFLOW_RETRY_MAX_ATTEMPTS=500

# Banking settings
BANKING_HOSTNAME="bank-api.example.com"
BANKING_TIMEOUT="10s"
```

## Testing the Compensation Path

To test the compensation logic (refund on failed deposit), you can modify the `Deposit` activity in `internal/workflow/activities.go` to use the `DepositThatFails` method instead of the normal `Deposit` method.

## Next Steps

- Add unit and integration tests
- Implement a real banking service client
- Add metrics and observability
- Create a web interface for initiating transfers

For more information on Temporal, please [read the documentation](https://docs.temporal.io/).
