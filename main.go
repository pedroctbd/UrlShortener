package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "shorturl.com/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"shorturl.com/handlers"
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
func addRoutes(r chi.Router) {

	r.Get("/swagger/*", httpSwagger.Handler())

	r.Post("/", handlers.CreateUrl())
	r.Get("/", handlers.RedirectUrl())

}

func run(ctx context.Context) error {

	r := chi.NewRouter()

	addRoutes(r)
	log.Println("Server running port 3000")
	return http.ListenAndServe(":3000", r)
}

func main() {

	ctx := context.Background()

	if err := run(ctx); err != nil {

		fmt.Println(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
