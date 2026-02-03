package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/cache"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/repository"
	hashutil "github.com/NugrahaPancaWibisana/solid-coffee-be/pkg/hash"
	jwtutil "github.com/NugrahaPancaWibisana/solid-coffee-be/pkg/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	authRepository *repository.AuthRepository
	redis          *redis.Client
	db             *pgxpool.Pool
}

func NewAuthService(authRepository *repository.AuthRepository, rdb *redis.Client, db *pgxpool.Pool) *AuthService {
	return &AuthService{authRepository: authRepository, redis: rdb, db: db}
}

func (as *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.User, error) {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, req.Email)
	if !matched {
		return dto.User{}, apperror.ErrInvalidEmailFormat
	}

	tx, err := as.db.Begin(ctx)
	if err != nil {
		log.Println(err.Error())
		return dto.User{}, err
	}
	defer tx.Rollback(ctx)

	data, err := as.authRepository.Login(ctx, tx, req)
	if err != nil {
		return dto.User{}, err
	}

	hasher := hashutil.Default()
	isValid, err := hasher.Verify(req.Password, data.Password)
	if err != nil {
		return dto.User{}, err
	}
	if !isValid {
		return dto.User{}, apperror.ErrInvalidCredential
	}

	err = as.authRepository.UpdateLastLogin(ctx, tx, data.ID)
	if err != nil {
		return dto.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("failed to commit", err.Error())
		return dto.User{}, err
	}

	res := dto.User{
		ID:        data.ID,
		Email:     data.Email,
		Role:      data.Role,
		LastLogin: nil,
	}

	if data.LastLoginAt.Valid {
		res.LastLogin = &data.LastLoginAt.Time
	}

	return res, nil
}

func (as *AuthService) GenerateJWT(ctx context.Context, user dto.User) (string, error) {
	claims := jwtutil.NewJWTClaims(user.ID, user.Role)
	return claims.GenToken()
}

func (as *AuthService) WhitelistToken(ctx context.Context, id int, token string) {
	cache.SetToken(ctx, as.redis, id, token)
}

func (as *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, req.Email)
	if !matched {
		return apperror.ErrInvalidEmailFormat
	}

	hasher := hashutil.Default()
	hashedPassword, err := hasher.Hash(req.Password)
	if err != nil {
		return err
	}

	req.Password = hashedPassword

	err = as.authRepository.Register(ctx, as.db, req)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) ForgotPassword(ctx context.Context, email string) error {
	err := as.authRepository.CheckEmailExists(ctx, as.db, email)
	if err != nil {
		return err
	}

	otp := make([]byte, 6)
	for i := range otp {
		otp[i] = byte('0' + rand.Intn(10))
	}

	rkey := fmt.Sprintf("%s:forgot-password:%s", os.Getenv("RDB_KEY"), string(otp))

	status := as.redis.Set(ctx, rkey, email, time.Minute*5)
	log.Println(string(otp))
	log.Println(status.Args()...)
	if status.Err() != nil {
		log.Println("caching failed:", status.Err())
	}

	return nil
}

func (as *AuthService) UpdatePassword(ctx context.Context, req dto.UpdateForgotPasswordRequest) error {
	rkey := fmt.Sprintf("%s:forgot-password:%s", os.Getenv("RDB_KEY"), req.Otp)

	email, err := as.redis.Get(ctx, rkey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return apperror.ErrOTPNotFound
		}
		return err
	}

	hasher := hashutil.Default()
	newHashedPassword, err := hasher.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	if err := as.authRepository.UpdatePassword(ctx, as.db, email, newHashedPassword); err != nil {
		return err
	}

	as.redis.Del(ctx, rkey)

	return nil
}
