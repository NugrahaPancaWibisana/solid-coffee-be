package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type configRDS struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

func InitRds() *redis.Client {
	config := configRDS{
		host:     os.Getenv("RDB_HOST"),
		port:     os.Getenv("RDB_PORT"),
		user:     os.Getenv("RDB_USERNAME"),
		password: os.Getenv("RDB_PASSWORD"),
		dbName:   os.Getenv("RDB_NAME"),
	}

	db, _ := strconv.Atoi(config.dbName)

	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.host, config.port),
		Username: config.user,
		Password: config.password,
		DB:       db,
	})
}
