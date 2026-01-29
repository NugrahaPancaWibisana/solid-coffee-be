package service

import (
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/repository"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	authRepository *repository.AuthRepository
	redis          *redis.Client
}

func NewAuthService(authRepository *repository.AuthRepository, rdb *redis.Client) *AuthService {
	return &AuthService{authRepository: authRepository, redis: rdb}
}
