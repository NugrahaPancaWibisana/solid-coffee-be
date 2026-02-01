package service

import (
	"context"
	"log"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/cache"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	userRepository *repository.UserRepository
	redis          *redis.Client
	db             *pgxpool.Pool
}

func NewUserService(userRepository *repository.UserRepository, rdb *redis.Client, db *pgxpool.Pool) *UserService {
	return &UserService{userRepository: userRepository, redis: rdb, db: db}
}

func (us *UserService) UpdateProfile(ctx context.Context, req dto.UpdateProfileRequest, path string, id int, token string) (string, error) {
	err := cache.CheckToken(ctx, us.redis, id, token)
	if err != nil {
		return "", err
	}

	tx, err := us.db.Begin(ctx)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	defer tx.Rollback(ctx)

	oldPath, err := us.userRepository.GetPhoto(ctx, tx, id)
	if err != nil {
		return "", err
	}

	if err := us.userRepository.UpdateProfile(ctx, tx, req, path, id); err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("failed to commit", err.Error())
		return "", err
	}

	return oldPath, nil
}
