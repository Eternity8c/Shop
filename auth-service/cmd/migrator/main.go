package main

import (
	// Библиотека для миграции

	// Драйвер для выполнения миграции SQLite 3
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	// Драйвер для получения миграции из файлов
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var dsn, migrationPath, migrationTable string

	flag.StringVar(&dsn, "dsn", "", "PostgreSQL connection string (DSN)")
	flag.StringVar(&migrationPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationTable, "migrations-table", "schema_migrations", "name of migrations table")
	flag.Parse()

	if dsn == "" {
		panic("dsn is required for PostgreSQL")
	}
	if migrationPath == "" {
		panic("migrations-path is required")
	}

	// Правильное добавление параметра
	var databaseURL string
	if strings.Contains(dsn, "?") {
		// Если уже есть параметры, добавляем через &
		databaseURL = fmt.Sprintf("%s&x-migrations-table=%s", dsn, migrationTable)
	} else {
		// Если нет параметров, добавляем через ?
		databaseURL = fmt.Sprintf("%s?x-migrations-table=%s", dsn, migrationTable)
	}

	m, err := migrate.New(
		"file://"+migrationPath,
		databaseURL,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}

		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
