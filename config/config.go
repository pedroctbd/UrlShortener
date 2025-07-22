package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Environment string

const (
	EnvTest       Environment = "test"
	EnvLocal      Environment = "local"
	EnvProduction Environment = "production"
)

type Config struct {
	Environment  Environment `required:"true" envconfig:"ENVIRONMENT" default:"local"`
	Development  bool        `required:"true" envconfig:"DEVELOPMENT"`
	App          App
	Postgres     Postgres
	JwtSecretKey string `required:"true" envconfig:"JWT_SECRET"`
}

type App struct {
	Name                    string        `required:"true" envconfig:"APP_NAME"`
	ID                      string        `required:"true" envconfig:"APP_ID"`
	GracefulShutdownTimeout time.Duration `required:"true" envconfig:"APP_GRACEFUL_SHUTDOWN_TIMEOUT"`
}

type Postgres struct {
	Host         string `envconfig:"DB_HOST"     default:"localhost"`
	User         string `envconfig:"DB_USER"     default:"postgres"`
	Password     string `envconfig:"DB_PASSWORD" default:"postgres"`
	DatabaseName string `envconfig:"DB_NAME"     default:"urlshortner"`
	Port         string `envconfig:"DB_PORT"     default:"5432"`
}

func New() (Config, error) {
	const operation = "Config.New"

	// 👇 Load .env file
	if err := godotenv.Load(); err != nil {
		// Optional: ignore if .env doesn't exist
		fmt.Println("No .env file found or failed to load")
	}

	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("%s -> %w", operation, err)
	}

	return cfg, nil
}
