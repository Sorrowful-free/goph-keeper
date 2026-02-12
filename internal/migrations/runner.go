package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

//go:embed sqlite/*.sql postgres/*.sql
var embedMigrations embed.FS

// RunUp выполняет миграции вверх (up). Использует отдельное подключение к БД для миграций,
// чтобы m.Close() не закрывал основное соединение GORM (иначе после миграций возникала бы ошибка "database is closed").
// dbType: "postgres" или "sqlite".
func RunUp(db *gorm.DB, dsn string, dbType string) error {
	// Проверяем, что основное подключение живо (миграции его не трогают)
	if _, err := db.DB(); err != nil {
		return fmt.Errorf("get underlying *sql.DB: %w", err)
	}

	isPostgres := dbType == "postgres"
	var databaseURL string
	var migrationsPath string

	if isPostgres {
		databaseURL = dsn
		if !strings.HasPrefix(databaseURL, "postgres://") && !strings.HasPrefix(databaseURL, "postgresql://") {
			databaseURL = "postgres://" + databaseURL
		}
		migrationsPath = "postgres"
	} else {
		if dsn == "" {
			dsn = "gophkeeper.db"
		}
		databaseURL = "sqlite3://" + dsn
		migrationsPath = "sqlite"
	}

	sourceDriver, err := iofs.New(embedMigrations, migrationsPath)
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	// Отдельное подключение только для миграций — migrate.Close() закроет его, не затрагивая GORM
	var migrateConn *sql.DB
	if isPostgres {
		migrateConn, err = sql.Open("postgres", databaseURL)
	} else {
		migrateConn, err = sql.Open("sqlite3", dsn)
	}
	if err != nil {
		return fmt.Errorf("open migration connection: %w", err)
	}
	// Не закрываем migrateConn вручную: migrate при m.Close() закрывает переданное WithInstance подключение

	var m *migrate.Migrate
	if isPostgres {
		driver, drvErr := postgres.WithInstance(migrateConn, &postgres.Config{})
		if drvErr != nil {
			return fmt.Errorf("create postgres driver: %w", drvErr)
		}
		m, err = migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	} else {
		driver, drvErr := sqlite3.WithInstance(migrateConn, &sqlite3.Config{})
		if drvErr != nil {
			return fmt.Errorf("create sqlite driver: %w", drvErr)
		}
		m, err = migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", driver)
	}
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations up: %w", err)
	}
	return nil
}
