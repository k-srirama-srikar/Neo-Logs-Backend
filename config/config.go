package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func ConnectDB() (*pgxpool.Pool, error) {
	LoadEnv()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	return pgxpool.New(nil, dsn)
}
