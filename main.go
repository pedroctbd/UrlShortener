package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"shorturl.com/config"
	_ "shorturl.com/docs"
	"shorturl.com/handlers"
	"shorturl.com/postgres"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
func addRoutes(r chi.Router, db *sql.DB) {

	r.Get("/swagger/*", httpSwagger.Handler())
	r.Post("/", handlers.CreateUrl(db))
	r.Get("/redirect", handlers.RedirectUrl())
	r.Get("/", handlers.ListExistingUrls(db))

}

func run(ctx context.Context, db *sql.DB) error {

	r := chi.NewRouter()

	addRoutes(r, db)
	log.Println("Server running port 3000")
	return http.ListenAndServe(":3000", r)
}

func main() {

	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load configurations: %v", err)
	}

	// Postgres
	postgresClient, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to start postgres: %v", err)
	}

	if err := run(ctx, postgresClient); err != nil {

		fmt.Println(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
