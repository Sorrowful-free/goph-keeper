@echo off
REM Скрипт для генерации Go кода из proto файлов (Windows)

set PROTO_DIR=proto
set OUT_DIR=proto

REM Проверяем наличие protoc
where protoc >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: protoc not found. Please install Protocol Buffers compiler.
    echo Download from: https://github.com/protocolbuffers/protobuf/releases
    exit /b 1
)

REM Проверяем наличие плагинов
where protoc-gen-go >nul 2>&1
if %errorlevel% neq 0 (
    echo Installing protoc-gen-go...
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
)

where protoc-gen-go-grpc >nul 2>&1
if %errorlevel% neq 0 (
    echo Installing protoc-gen-go-grpc...
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
)

REM Генерируем код
echo Generating Go code from proto files...
protoc --go_out=%OUT_DIR% --go_opt=paths=source_relative ^
    --go-grpc_out=%OUT_DIR% --go-grpc_opt=paths=source_relative ^
    %PROTO_DIR%\*.proto

echo Proto files generated successfully!
