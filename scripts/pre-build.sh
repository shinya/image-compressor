#!/bin/bash
set -e

echo "Running pre-build script..."

# 依存関係の確認
echo "Checking dependencies..."
go mod verify

# テストの実行
echo "Running tests..."
go test ./...

# 静的解析
echo "Running static analysis..."
go vet ./...

# コードフォーマットの確認
echo "Checking code formatting..."
if [ -n "$(gofmt -l .)" ]; then
    echo "Code is not formatted. Please run 'go fmt ./...'"
    exit 1
fi

echo "Pre-build script completed successfully!" 