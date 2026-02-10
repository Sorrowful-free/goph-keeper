.PHONY: build-server build-client proto clean test mocks migrate-up migrate-down migrate-create

# Версия и дата сборки
VERSION ?= 1.0.0
BUILD_DATE ?= $(shell date +%Y-%m-%d)

# Протобуфы
PROTO_DIR = proto
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)
PROTO_GO_FILES = $(PROTO_FILES:.proto=.pb.go)

# Сборка сервера
build-server:
	@echo "Building server..."
	@mkdir -p bin
	go build -o bin/server ./cmd/server

# Сборка клиента с версией
build-client:
	@echo "Building client..."
	@mkdir -p bin
	go build -ldflags "-X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)" -o bin/client ./cmd/client

# Генерация протобуфов
proto:
	@echo "Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

# Очистка
clean:
	rm -rf bin/
	rm -f $(PROTO_DIR)/*.pb.go

# Тесты
test:
	go test ./...

# Генерация моков (go.uber.org/mock)
mocks:
	@echo "Generating mocks..."
	go generate ./internal/domain/repository/...

# Установка зависимостей
deps:
	go mod download
	go mod tidy

# Миграции (go-migrate). DSN по умолчанию: sqlite3://gophkeeper.db
# Установка CLI: go install -tags 'no_postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
MIGRATE_DSN ?= sqlite3://gophkeeper.db
MIGRATE_PATH_SQLITE = internal/migrations/sqlite
MIGRATE_PATH_POSTGRES = internal/migrations/postgres

migrate-up-sqlite:
	migrate -path $(MIGRATE_PATH_SQLITE) -database "$(MIGRATE_DSN)" up

migrate-up-postgres:
	migrate -path $(MIGRATE_PATH_POSTGRES) -database "$(MIGRATE_DSN)" up

migrate-down-sqlite:
	migrate -path $(MIGRATE_PATH_SQLITE) -database "$(MIGRATE_DSN)" down

migrate-down-postgres:
	migrate -path $(MIGRATE_PATH_POSTGRES) -database "$(MIGRATE_DSN)" down

# Создать новую миграцию (имя передать: make migrate-create name=add_foo)
migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=description"; exit 1)
	@mkdir -p $(MIGRATE_PATH_SQLITE) $(MIGRATE_PATH_POSTGRES)
	@n=$$(ls $(MIGRATE_PATH_SQLITE) 2>/dev/null | grep -E '^[0-9]+_' | sed 's/_.*//' | sort -n | tail -1); n=$${n:-0}; n=$$(($$n+1)); \
	id=$$(printf '%06d' $$n); \
	touch $(MIGRATE_PATH_SQLITE)/$${id}_$(name).up.sql $(MIGRATE_PATH_SQLITE)/$${id}_$(name).down.sql; \
	touch $(MIGRATE_PATH_POSTGRES)/$${id}_$(name).up.sql $(MIGRATE_PATH_POSTGRES)/$${id}_$(name).down.sql; \
	echo "Created migrations: $${id}_$(name)"
