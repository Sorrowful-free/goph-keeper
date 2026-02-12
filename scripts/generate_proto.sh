#!/bin/bash

# Скрипт для генерации Go кода из proto файлов

PROTO_DIR="proto"
OUT_DIR="proto"

# Проверяем наличие protoc
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc not found. Please install Protocol Buffers compiler."
    echo "On macOS: brew install protobuf"
    echo "On Ubuntu: sudo apt-get install protobuf-compiler"
    exit 1
fi

# Проверяем наличие плагинов
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Генерируем код
echo "Generating Go code from proto files..."
protoc --go_out=$OUT_DIR --go_opt=paths=source_relative \
    --go-grpc_out=$OUT_DIR --go-grpc_opt=paths=source_relative \
    $PROTO_DIR/*.proto

echo "Proto files generated successfully!"
