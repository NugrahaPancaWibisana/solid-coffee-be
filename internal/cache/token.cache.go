package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/redis/go-redis/v9"
)

var (
	rkey = "auth:token:"
	exp  = 24 * time.Hour
)

func CheckToken(ctx context.Context, rdb *redis.Client, id int, token string) error {
	key := fmt.Sprintf("%s:%s%d", os.Getenv("RDB_KEY"),rkey, id)
	tokenCache, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return apperror.ErrSessionExpired
	}

	if err != nil {
		log.Println("Redis error:", err.Error())
		return apperror.ErrInternal
	}

	if tokenCache != token {
		return apperror.ErrInvalidSession
	}

	return nil
}

func SetToken(ctx context.Context, rdb *redis.Client, id int, token string) {
	key := fmt.Sprintf("%s:%s%d", os.Getenv("RDB_KEY"),rkey, id)

	status := rdb.Set(ctx, key, token, exp)
	if status.Err() != nil {
		log.Println("caching failed:", status.Err())
	}

}

func DeleteToken(ctx context.Context, rdb *redis.Client, id int) error {
	key := fmt.Sprintf("%s:%s%d", os.Getenv("RDB_KEY"),rkey, id)

	_, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return apperror.ErrLogoutFailed
	}

	if err != nil {
		log.Println("Redis error:", err.Error())
		return apperror.ErrInternal
	}

	err = rdb.Del(ctx, key).Err()

	if err != nil {
		return apperror.ErrLogoutFailed
	}

	return nil
}
