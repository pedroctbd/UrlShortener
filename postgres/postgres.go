package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"shorturl.com/config"
)

// New connects to the Postgres database and performs migrations.
func New(ctx context.Context, config config.Postgres) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.DatabaseName,
	)

	fmt.Println(connString)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, err
}
