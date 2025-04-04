#!/bin/sh
set -e

# 配置参数
PROJECT_DIR="/home/mystks_api"         # 代码所在目录
OUTPUT_DIR="/app/stonks-api"           # 编译输出目录
BINARY_NAME="stonks-api"               # 生成的二进制文件名称

echo "切换到项目目录：$PROJECT_DIR"
cd "$PROJECT_DIR"

echo "拉取最新代码..."
git pull

echo "开始编译 Go 程序..."
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64

# 编译项目并将二进制文件输出到指定目录
go build -o "$OUTPUT_DIR/$BINARY_NAME" main.go
echo "编译成功，二进制文件位于 $OUTPUT_DIR/$BINARY_NAME"