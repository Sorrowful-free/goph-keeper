# Настройка Proto файлов

## Генерация Go кода из Proto файлов

После генерации proto файлов необходимо обновить `cmd/server/main.go` для регистрации сервисов.

### Шаги:

1. **Сгенерируйте proto файлы:**

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

2. **Обновите `cmd/server/main.go`:**

   Добавьте импорт:
   ```go
   import (
       // ... другие импорты
       "github.com/gophkeeper/gophkeeper/proto"
   )
   ```

   Раскомментируйте регистрацию сервисов:
   ```go
   // Регистрируем сервисы
   authService := server.NewAuthService(st)
   dataService := server.NewDataService(st)

   proto.RegisterAuthServiceServer(grpcServer, authService)
   proto.RegisterDataServiceServer(grpcServer, dataService)
   ```

3. **Проверьте, что файлы сгенерированы:**

   После генерации должны появиться файлы:
   - `proto/gophkeeper.pb.go`
   - `proto/gophkeeper_grpc.pb.go`

## Структура сгенерированных файлов

После генерации структура будет следующей:

```
proto/
├── gophkeeper.proto          # Исходный proto файл
├── gophkeeper.pb.go          # Сгенерированные структуры
└── gophkeeper_grpc.pb.go     # Сгенерированные gRPC сервисы
```

## Проверка генерации

Убедитесь, что в сгенерированных файлах есть:
- `RegisterAuthServiceServer` функция
- `RegisterDataServiceServer` функция
- Все типы данных из proto файла

## Устранение проблем

### Ошибка: "protoc: command not found"

Установите Protocol Buffers Compiler:
- **macOS**: `brew install protobuf`
- **Ubuntu**: `sudo apt-get install protobuf-compiler`
- **Windows**: Скачайте с [GitHub Releases](https://github.com/protocolbuffers/protobuf/releases)

### Ошибка: "protoc-gen-go: program not found"

Установите плагины:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Убедитесь, что `$GOPATH/bin` или `$HOME/go/bin` в PATH.

### Ошибки компиляции после генерации

1. Убедитесь, что все зависимости установлены: `go mod tidy`
2. Проверьте, что импорты в `cmd/server/main.go` правильные
3. Убедитесь, что сервисы правильно реализуют интерфейсы из proto
