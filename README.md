# GophKeeper

Менеджер паролей GophKeeper - клиент-серверная система для безопасного хранения приватной информации.

## Возможности

- Регистрация и аутентификация пользователей
- Хранение различных типов данных:
  - Пары логин/пароль
  - Произвольные текстовые данные
  - Произвольные бинарные данные
  - Данные банковских карт
- Синхронизация данных между устройствами
- TUI интерфейс для клиента
- Безопасное шифрование данных

## Технологии

- Go 1.21+
- gRPC для взаимодействия клиента и сервера
- PostgreSQL/SQLite для хранения данных
- Bubbletea для TUI интерфейса

## Быстрый старт

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

### 2. Установка зависимостей

```bash
go mod download
go mod tidy
```

### 3. Сборка

**Сервер:**
```bash
go build -o bin/server ./cmd/server
```

**Клиент:**
```bash
go build -ldflags "-X main.version=1.0.0 -X main.buildDate=$(date +%Y-%m-%d)" -o bin/client ./cmd/client
```

Или используйте Makefile:
```bash
make build-server
make build-client VERSION=1.0.0
```

Подробные инструкции см. в [BUILD.md](BUILD.md)

## Использование

### Запуск сервера

```bash
./bin/server --config config.yaml
```

### Запуск клиента

```bash
./bin/client --server localhost:50051
```

## Структура проекта

```
.
├── cmd/
│   ├── server/     # Серверное приложение
│   └── client/     # Клиентское приложение
├── internal/
│   ├── server/     # Серверная логика
│   ├── client/     # Клиентская логика
│   ├── models/     # Модели данных
│   ├── storage/    # Работа с БД
│   └── crypto/     # Криптография
├── proto/          # gRPC протобуфы
└── pkg/            # Общие пакеты

```
