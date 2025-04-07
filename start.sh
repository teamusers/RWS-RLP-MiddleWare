#!/bin/bash

# Configure environment variables
export DALINK_GO_CONFIG_PATH=/app/stonks-api/prod.yml

# Process name
PROCESS_NAME="./rlp-middleware-api"
LOG_FILE="output.log"

# Check if the process exists
PID=$(pgrep -f $PROCESS_NAME)

if [ -n "$PID" ]; then
  echo "Killing existing process $PROCESS_NAME with PID $PID"
  kill -9 $PID
  sleep 2
else
  echo "No existing process found for $PROCESS_NAME"
fi

# Start a new process
echo "Starting new process $PROCESS_NAME"
nohup ./$PROCESS_NAME > $LOG_FILE 2>&1 &

# Get the new process ID
NEW_PID=$(pgrep -f $PROCESS_NAME)
if [ -n "$NEW_PID" ]; then
  echo "New process $PROCESS_NAME started with PID $NEW_PID"
else
  echo "Failed to start new process $PROCESS_NAME"
fi
