# Инструкции по сборке GophKeeper

## Предварительные требования

1. **Go 1.21 или выше**
   - Скачайте с [golang.org](https://golang.org/dl/)

2. **Protocol Buffers Compiler (protoc)**
   - **macOS**: `brew install protobuf`
   - **Ubuntu/Debian**: `sudo apt-get install protobuf-compiler`
   - **Windows**: Скачайте с [GitHub Releases](https://github.com/protocolbuffers/protobuf/releases)

3. **Go плагины для protoc**
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

## Сборка проекта

### 1. Генерация proto файлов

**Linux/macOS:**
```bash
chmod +x scripts/generate_proto.sh
./scripts/generate_proto.sh
```

**Windows:**
```cmd
scripts\generate_proto.bat
```

Или вручную:
```bash
protoc --go_out=proto --go_opt=paths=source_relative \
    --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
    proto/*.proto
```

### 2. Установка зависимостей

```bash
go mod download
go mod tidy
```

### 3. Сборка сервера

```bash
go build -o bin/server ./cmd/server
```

### 4. Сборка клиента

**С версией и датой сборки:**
```bash
go build -ldflags "-X main.version=1.0.0 -X main.buildDate=$(date +%Y-%m-%d)" -o bin/client ./cmd/client
```

**Windows:**
```cmd
go build -ldflags "-X main.version=1.0.0 -X main.buildDate=%date%" -o bin\client.exe .\cmd\client
```

**Или используйте Makefile:**
```bash
make build-server
make build-client VERSION=1.0.0
```

## Запуск

### Сервер

```bash
# С SQLite (по умолчанию)
./bin/server

# С PostgreSQL
./bin/server --dsn "host=localhost user=postgres password=postgres dbname=gophkeeper sslmode=disable"

# На другом порту
./bin/server --port 8080
```

### Клиент

```bash
# Подключение к серверу
./bin/client --server localhost:50051

# Просмотр версии
./bin/client --version
```

## Использование

1. Запустите сервер
2. Запустите клиент
3. В клиенте:
   - Зарегистрируйтесь или войдите
   - Используйте меню для работы с данными
   - Навигация: стрелки ↑↓, Enter для выбора, Esc для возврата, q для выхода

## Структура проекта

```
.
├── cmd/
│   ├── server/          # Серверное приложение
│   └── client/          # Клиентское приложение
├── internal/
│   ├── server/          # Серверная логика (gRPC сервисы)
│   ├── client/          # Клиентская логика
│   │   └── tui/         # TUI интерфейс
│   ├── models/          # Модели данных
│   ├── storage/         # Работа с БД
│   └── crypto/          # Криптография
├── proto/               # gRPC протобуфы
└── scripts/             # Скрипты сборки
```

## Переменные окружения

- `DB_TYPE` - тип БД: `postgres` или `sqlite` (по умолчанию)
- `JWT_SECRET` - секретный ключ для JWT (в продакшене обязательно!)

## Примечания

- По умолчанию используется SQLite для простоты разработки
- В продакшене обязательно используйте PostgreSQL и установите `JWT_SECRET`
- Данные шифруются на клиенте перед отправкой на сервер
- Сервер хранит только зашифрованные данные
