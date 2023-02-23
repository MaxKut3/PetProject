package config

import (
	"log"
	"os"

	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
)

type Config struct {
	Conn pgx.ConnConfig
}

func NewConfig() *Config {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf(".env файл не прочитан: %e", err)
	}

	return &Config{
		Conn: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5434,
			Database: "postgres",
			User:     "postgres",
			Password: os.Getenv("DB"),
		},
	}
}
