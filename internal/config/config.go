package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Config — конфигурация сервера (всё, что парсится из флагов и переменных окружения).
type Config struct {
	// Server
	Port     string // порт (флаг -port)
	GrpcAddr string // адрес gRPC (флаг -addr, переопределяет port)
	Address  string // итоговый адрес слушателя, например ":50051"

	// Database
	DSN     string // строка подключения к БД (флаг -dsn)
	DBType  string // "postgres" или "sqlite"
	DefaultDSN string // DSN по умолчанию, если не задан (sqlite: gophkeeper.db)

	// Security (JWT)
	JWTSecret          []byte        // секрет для подписи JWT (env JWT_SECRET)
	AccessTokenExpiry  time.Duration // время жизни access токена (env ACCESS_TOKEN_EXPIRY)
	RefreshTokenExpiry time.Duration // время жизни refresh токена (env REFRESH_TOKEN_EXPIRY)
}

const (
	DBTypePostgres = "postgres"
	DBTypeSQLite   = "sqlite"
	defaultPort    = "50051"
	defaultDSN     = "gophkeeper.db"
	defaultJWT     = "your-secret-key-change-in-production"
	defaultAccess  = 15 * time.Minute
	defaultRefresh = 7 * 24 * time.Hour
)

// Load парсит флаги и переменные окружения, заполняет и возвращает Config.
// Флаги: -port, -dsn, -addr.
// Env: DB_TYPE, JWT_SECRET, ACCESS_TOKEN_EXPIRY, REFRESH_TOKEN_EXPIRY.
func Load() *Config {
	port := flag.String("port", defaultPort, "Server port")
	dsn := flag.String("dsn", "", "Database connection string (default: SQLite)")
	grpcAddr := flag.String("addr", "", "gRPC server address (overrides port)")
	flag.Parse()

	cfg := &Config{
		Port:       *port,
		GrpcAddr:   *grpcAddr,
		DSN:        *dsn,
		DefaultDSN: defaultDSN,
	}

	// Итоговый адрес
	if cfg.GrpcAddr != "" {
		cfg.Address = cfg.GrpcAddr
	} else {
		cfg.Address = fmt.Sprintf(":%s", cfg.Port)
	}

	// Тип БД: env DB_TYPE или по префиксу DSN
	if os.Getenv("DB_TYPE") == DBTypePostgres {
		cfg.DBType = DBTypePostgres
	} else if cfg.DSN != "" && len(cfg.DSN) >= 4 && strings.HasPrefix(cfg.DSN, "post") {
		cfg.DBType = DBTypePostgres
	} else {
		cfg.DBType = DBTypeSQLite
	}

	if cfg.DSN == "" {
		cfg.DSN = cfg.DefaultDSN
	}

	// JWT: из env или дефолты
	if s := os.Getenv("JWT_SECRET"); s != "" {
		cfg.JWTSecret = []byte(s)
	} else {
		cfg.JWTSecret = []byte(defaultJWT)
	}
	if s := os.Getenv("ACCESS_TOKEN_EXPIRY"); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			cfg.AccessTokenExpiry = d
		} else {
			cfg.AccessTokenExpiry = defaultAccess
		}
	} else {
		cfg.AccessTokenExpiry = defaultAccess
	}
	if s := os.Getenv("REFRESH_TOKEN_EXPIRY"); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			cfg.RefreshTokenExpiry = d
		} else {
			cfg.RefreshTokenExpiry = defaultRefresh
		}
	} else {
		cfg.RefreshTokenExpiry = defaultRefresh
	}

	return cfg
}
