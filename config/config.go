package config

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LoadEnv loads environment variables from a .env file.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found.")
	}
}

// ConnectDB creates and returns a PostgreSQL connection pool.
func ConnectDB() (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	return pgxpool.New(context.Background(), dsn)
}
