package repository

import (
	"context"
	"errors"
	"log"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderRepo interface {
	Login(ctx context.Context, req dto.LoginRequest, db DBTX) (model.User, error)
	UpdateLastLogin(ctx context.Context, db DBTX, id int) error
}

type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type AuthRepository struct {
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (ar *AuthRepository) Login(ctx context.Context, db DBTX, req dto.LoginRequest) (model.User, error) {
	query := `
		SELECT
		    id,
		    email,
		    password,
		    role,
		    lastlogin_at
		FROM
		    users u
		WHERE
		    email = $1;
	`

	row := db.QueryRow(ctx, query, req.Email)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.LastLoginAt,
	)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, apperror.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (ar *AuthRepository) UpdateLastLogin(ctx context.Context, db DBTX, id int) error {
	query := `
		UPDATE users
		SET
		    lastlogin_at = NOW()
		WHERE
		    id = $1;
	`

	_, err := db.Exec(ctx, query, id)
	if err != nil {
		log.Println(err.Error())
		return apperror.ErrUpdateLastLogin
	}

	return nil
}
