package config

import (
	"flag"
	"os"
)

// ClientConfig — конфигурация клиента (всё, что парсится из флагов и переменных окружения).
type ClientConfig struct {
	// Server — адрес gRPC-сервера (флаг -server или env SERVER_ADDRESS).
	Server string
}

const defaultServer = "localhost:50051"

// LoadClient парсит флаги и переменные окружения, заполняет и возвращает ClientConfig.
// Флаг: -server.
// Env: SERVER_ADDRESS (переопределяет флаг).
func LoadClient() *ClientConfig {
	server := flag.String("server", defaultServer, "Server address")
	flag.Parse()

	cfg := &ClientConfig{
		Server: *server,
	}
	if s := os.Getenv("SERVER_ADDRESS"); s != "" {
		cfg.Server = s
	}
	return cfg
}
