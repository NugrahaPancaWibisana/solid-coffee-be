package service

import (
	"context"
	"log"
	"regexp"

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

	if e := tx.Commit(ctx); e != nil {
		log.Println("failed to commit", e.Error())
		return dto.User{}, e
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
