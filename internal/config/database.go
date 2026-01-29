package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type configDB struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

func InitDB() (*pgxpool.Pool, error) {
	config := configDB{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		user:     os.Getenv("DB_USERNAME"),
		password: os.Getenv("DB_PASSWORD"),
		dbName:   os.Getenv("DB_NAME"),
	}

	return pgxpool.New(
		context.Background(),
		fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			config.user,
			config.password,
			config.host,
			config.port,
			config.dbName,
		),
	)
}
