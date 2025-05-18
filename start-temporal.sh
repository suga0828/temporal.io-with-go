#!/bin/bash
#==============================================================================
# Temporal.io Money Transfer Application Setup Script
#==============================================================================

#------------------------------------------------------------------------------
# Load environment variables from .env file if it exists
#------------------------------------------------------------------------------
ENV_FILE=".env"
if [ -f "$ENV_FILE" ]; then
  echo "Loading environment variables from $ENV_FILE"
  source "$ENV_FILE"
else
  echo "Warning: $ENV_FILE file not found. Using default values."
  # Copy the example file if it doesn't exist
  if [ -f ".env.example" ] && [ ! -f "$ENV_FILE" ]; then
    echo "Creating $ENV_FILE from .env.example"
    cp .env.example "$ENV_FILE"
    source "$ENV_FILE"
  fi
fi

# Set defaults for required variables if not set in .env
export TEMPORAL_UI_PORT=${TEMPORAL_UI_PORT:-8080}
export TEMPORAL_PORT=${TEMPORAL_PORT:-7233}
export TEMPORAL_DB_FILE=${TEMPORAL_DB_FILE:-"temporal.db"}
export TEMPORAL_NAMESPACE=${TEMPORAL_NAMESPACE:-"default"}
export WORKFLOW_ID=${WORKFLOW_ID:-"money-transfer-workflow"}
export TEMPORAL_ADDRESS=${TEMPORAL_ADDRESS:-"localhost:${TEMPORAL_PORT}"}
export WORKFLOW_TASK_QUEUE=${WORKFLOW_TASK_QUEUE:-"TRANSFER_MONEY_TASK_QUEUE"}
export LOG_MODE=${LOG_MODE:-"development"}


#------------------------------------------------------------------------------
# Helper Functions
#------------------------------------------------------------------------------
# Function to check if a process is running on a specific port
check_port() {
  if lsof -i :$1 > /dev/null; then
    return 0  # Port is in use
  else
    return 1  # Port is free
  fi
}


#------------------------------------------------------------------------------
# Component Management Functions
#------------------------------------------------------------------------------
# Function to start the Temporal server
start_temporal_server() {
  echo "Starting Temporal server on UI port $TEMPORAL_UI_PORT and service port $TEMPORAL_PORT..."
  
  if check_port $TEMPORAL_UI_PORT; then
    echo "Port $TEMPORAL_UI_PORT is already in use. Temporal UI might be running already."
  elif check_port $TEMPORAL_PORT; then
    echo "Port $TEMPORAL_PORT is already in use. Temporal service might be running already."
  else
    # Start Temporal server in the background
    temporal server start-dev --db-filename $TEMPORAL_DB_FILE --ui-port $TEMPORAL_UI_PORT --port $TEMPORAL_PORT &
    TEMPORAL_PID=$!
    echo "Temporal server started with PID: $TEMPORAL_PID"
    
    # Wait for server to be ready
    echo "Waiting for Temporal server to be ready..."
    sleep 5
    
    # Verify server is actually running
    if ! ps -p $TEMPORAL_PID > /dev/null; then
      echo "ERROR: Temporal server failed to start or crashed!"
      exit 1
    fi
    
    # Check if server is accepting connections
    for i in {1..5}; do
      if nc -z localhost $TEMPORAL_PORT 2>/dev/null; then
        echo "Temporal server is ready and accepting connections."
        return 0
      fi
      echo "Waiting for Temporal server to accept connections (attempt $i/5)..."
      sleep 2
    done
    
    echo "ERROR: Temporal server is not accepting connections after multiple attempts."
    echo "Check server logs for more information."
    kill $TEMPORAL_PID 2>/dev/null || true
    exit 1
  fi
}


# Function to start the worker
start_worker() {
  # Check if Temporal server is accepting connections
  if ! nc -z localhost $TEMPORAL_PORT 2>/dev/null; then
    echo "ERROR: Cannot connect to Temporal server at localhost:$TEMPORAL_PORT"
    echo "Worker cannot be started without a running Temporal server."
    return 1
  fi

  echo "Starting Temporal worker..."
  cd cmd/worker && go run main.go &
  WORKER_PID=$!
  echo "Worker started with PID: $WORKER_PID"
  
  # Verify worker is actually running
  sleep 2
  if ! ps -p $WORKER_PID > /dev/null; then
    echo "ERROR: Worker failed to start or crashed!"
    return 1
  fi
  
  return 0
}


# Function to execute the workflow
execute_workflow() {
  # Check if Temporal server is accepting connections
  if ! nc -z localhost $TEMPORAL_PORT 2>/dev/null; then
    echo "ERROR: Cannot connect to Temporal server at localhost:$TEMPORAL_PORT"
    echo "Workflow cannot be executed without a running Temporal server."
    return 1
  fi

  echo "Executing money transfer workflow with ID: $WORKFLOW_ID..."
  cd cmd/starter && go run main.go
  WORKFLOW_STATUS=$?
  
  if [ $WORKFLOW_STATUS -eq 0 ]; then
    echo "Workflow executed successfully!"
  else
    echo "Workflow execution failed with status: $WORKFLOW_STATUS"
  fi
  
  return $WORKFLOW_STATUS
}


# Function to clean up processes
cleanup() {
  echo ""
  echo "-------------------------------------"
  echo "  CLEANING UP PROCESSES"
  echo "-------------------------------------"
  
  if [ ! -z "$WORKER_PID" ] && ps -p $WORKER_PID > /dev/null; then
    echo "Stopping worker (PID: $WORKER_PID)..."
    kill -9 $WORKER_PID 2>/dev/null || true
    wait $WORKER_PID 2>/dev/null || true
  fi
  
  if [ ! -z "$TEMPORAL_PID" ] && ps -p $TEMPORAL_PID > /dev/null; then
    echo "Stopping Temporal server (PID: $TEMPORAL_PID)..."
    kill -9 $TEMPORAL_PID 2>/dev/null || true
    wait $TEMPORAL_PID 2>/dev/null || true
  fi
  
  echo "Cleanup complete."
  exit 0
}


#------------------------------------------------------------------------------
# Main Script Execution
#------------------------------------------------------------------------------
# Ensure cleanup happens on script exit
trap cleanup EXIT INT TERM

# Print configuration information
echo "=== Temporal Money Transfer Application Setup ===="
echo "Temporal Service Port: $TEMPORAL_PORT"
echo "UI Port: $TEMPORAL_UI_PORT"
echo "Database File: $TEMPORAL_DB_FILE"
echo "Namespace: $TEMPORAL_NAMESPACE"
echo "Workflow ID: $WORKFLOW_ID"
echo "Task Queue: $WORKFLOW_TASK_QUEUE"
echo "=================================================="

# Start components
echo ""
echo "-------------------------------------"
echo "  STARTING TEMPORAL SERVER"
echo "-------------------------------------"
start_temporal_server || exit 1

echo ""
echo "-------------------------------------"
echo "  STARTING WORKER"
echo "-------------------------------------"
start_worker || echo "WARNING: Worker failed to start properly."

echo ""
echo "-------------------------------------"
echo "  EXECUTING WORKFLOW"
echo "-------------------------------------"
execute_workflow

echo ""
echo "-------------------------------------"
echo "  MONITORING MODE"
echo "-------------------------------------"
# Keep the script running to maintain the background processes
echo "Press Ctrl+C to stop all processes and exit"
while true; do
  # Check if processes are still running
  if [ ! -z "$TEMPORAL_PID" ] && ! ps -p $TEMPORAL_PID > /dev/null; then
    echo "ERROR: Temporal server process has stopped unexpectedly!"
    break
  fi
  
  if [ ! -z "$WORKER_PID" ] && ! ps -p $WORKER_PID > /dev/null; then
    echo "WARNING: Worker process has stopped unexpectedly."
  fi
  
  sleep 5
done

# If we get here, something went wrong
echo "Exiting due to unexpected process termination."
cleanup
