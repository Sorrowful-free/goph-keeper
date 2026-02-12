# Миграции (golang-migrate)

Миграции хранятся в `internal/migrations/sqlite` и `internal/migrations/postgres` и вшиваются в бинарник (embed). При старте сервера миграции **выполняются автоматически** (все pending up).

## Запуск миграций из кода

При старте `cmd/server` вызывается `migrations.RunUp(db, dsn)` после подключения к БД. Тип БД определяется по DSN или переменной окружения `DB_TYPE=postgres`.

## CLI (migrate)

Для ручного запуска миграций установите [migrate](https://github.com/golang-migrate/migrate):

```bash
# С поддержкой SQLite (требуется CGO)
go install -tags 'no_postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Или с поддержкой PostgreSQL
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### SQLite

```bash
# В корне проекта
migrate -path internal/migrations/sqlite -database "sqlite3://gophkeeper.db" up
migrate -path internal/migrations/sqlite -database "sqlite3://gophkeeper.db" down
```

Через Makefile (по умолчанию DSN = sqlite3://gophkeeper.db):

```bash
make migrate-up-sqlite
make migrate-down-sqlite
```

### PostgreSQL

```bash
migrate -path internal/migrations/postgres -database "postgres://user:pass@localhost:5432/gophkeeper?sslmode=disable" up
make migrate-up-postgres MIGRATE_DSN="postgres://user:pass@localhost:5432/gophkeeper?sslmode=disable"
```

## Создание новой миграции

Добавьте пары файлов в оба каталога (sqlite и postgres) с одинаковым номером и именем:

- `internal/migrations/sqlite/000003_описание.up.sql`
- `internal/migrations/sqlite/000003_описание.down.sql`
- `internal/migrations/postgres/000003_описание.up.sql`
- `internal/migrations/postgres/000003_описание.down.sql`

На Unix можно использовать:

```bash
make migrate-create name=add_some_table
```

После этого отредактируйте созданные `.sql` файлы.
