package repository

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/model"
	"github.com/jackc/pgx/v5"
)

type AuthRepo interface {
	Login(ctx context.Context, req dto.LoginRequest, db DBTX) (model.User, error)
	UpdateLastLogin(ctx context.Context, db DBTX, id int) error
}

type AuthRepository struct{}

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

func (ar *AuthRepository) Register(ctx context.Context, db DBTX, req dto.RegisterRequest) error {
	query := `
		INSERT INTO
		    public.users (fullname, email, password)
		VALUES
		    ($1, $2, $3)
	`

	_, err := db.Exec(ctx, query, req.Fullname, req.Email, req.Password)
	if err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "duplicate") {
			return apperror.ErrEmailAlreadyExists
		}
		return apperror.ErrRegisterUser
	}

	return nil
}

func (ar *AuthRepository) CheckEmailExists(ctx context.Context, db DBTX, email string) error {
	query := "SELECT email FROM users WHERE email = $1 AND deleted_at IS NULL"

	var foundEmail string
	err := db.QueryRow(ctx, query, email).Scan(&foundEmail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrUserNotFound
		}
		return err
	}

	return nil
}

func (ar *AuthRepository) UpdatePassword(ctx context.Context, db DBTX, email, password string) error {
	query := `
		UPDATE users
		SET password = $1, updated_at = NOW()
		WHERE email = $2 AND deleted_at IS NULL
	`

	ct, err := db.Exec(ctx, query, password, email)
	if err != nil {
		log.Println(err.Error())
		return apperror.ErrUpdatePassword
	}

	if ct.RowsAffected() == 0 {
		return apperror.ErrUserNotFound
	}

	return nil
}
