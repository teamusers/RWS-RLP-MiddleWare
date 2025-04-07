#!/bin/sh
set -e
REMOTE_USER="root"
REMOTE_HOST="47.84.41.204"
REMOTE_DIR="/app/rlp-middleware-api"
echo "Start executing update and build......"
ssh ${REMOTE_USER}@${REMOTE_HOST} "cd ${REMOTE_DIR} && ./build.sh"
echo "Start launching the application......"
ssh ${REMOTE_USER}@${REMOTE_HOST} "cd ${REMOTE_DIR} && ./start.sh"
echo "Update and launch operation completed!"