package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"shorturl.com/config"
)

// New connects to the Postgres database and performs migrations.
func New(ctx context.Context, config config.Postgres) (*sql.DB, error) {

	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", config.User, config.Password, config.Host, config.Port, config.DatabaseName)

	fmt.Println("Connecting with:", connString)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %w", err)
	}

	// Check DB connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping error: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate driver error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://postgres/migrations",
		"postgres", driver)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("migration init error: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		db.Close()
		return nil, fmt.Errorf("migration error: %w", err)
	}

	return db, nil
}
