#!/bin/sh
set -e
REMOTE_USER="root"
REMOTE_HOST="47.84.41.204"
REMOTE_DIR="/app/stonks-api"
echo "开始执行更新并构建......"
ssh ${REMOTE_USER}@${REMOTE_HOST} "cd ${REMOTE_DIR} && ./build.sh"
echo "开始启动应用......"
ssh ${REMOTE_USER}@${REMOTE_HOST} "cd ${REMOTE_DIR} && ./start.sh"
echo "更新并启动操作完成。"