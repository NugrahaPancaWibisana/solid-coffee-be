package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/model"
	"github.com/jackc/pgx/v5"
)

type UserRepo interface {
	GetPhoto(ctx context.Context, db DBTX, id int) (string, error)
	UpdateProfile(ctx context.Context, db DBTX, req dto.UpdateProfileRequest, path string, id int) error
	GetPasswordByUserID(ctx context.Context, db DBTX, id int) (string, error)
	UpdatePassword(ctx context.Context, db DBTX, id int, password string) error
}

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) GetPhoto(ctx context.Context, db DBTX, id int) (string, error) {
	query := "SELECT photo FROM users WHERE id = $1;"

	row := db.QueryRow(ctx, query, id)

	var photo *string
	err := row.Scan(&photo)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	if photo == nil {
		return "", nil
	}

	return *photo, nil
}

func (ur *UserRepository) UpdateProfile(ctx context.Context, db DBTX, req dto.UpdateProfileRequest, path string, id int) error {
	var sb strings.Builder
	sb.WriteString("UPDATE users SET ")
	args := []any{}

	if path != "" {
		if len(args) > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "photo = $%d", len(args)+1)
		args = append(args, path)
	}

	if req.Fullname != "" {
		if len(args) > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "fullname = $%d", len(args)+1)
		args = append(args, req.Fullname)
	}

	if req.Phone != "" {
		if len(args) > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "phone = $%d", len(args)+1)
		args = append(args, req.Phone)
	}

	if req.Address != "" {
		if len(args) > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "address = $%d", len(args)+1)
		args = append(args, req.Address)
	}

	if len(args) == 0 {
		return apperror.ErrNoFieldsToUpdate
	}

	fmt.Fprintf(&sb, " WHERE id = $%d", len(args)+1)
	args = append(args, id)

	_, err := db.Exec(ctx, sb.String(), args...)
	if err != nil {
		log.Println(err.Error())
		return apperror.ErrUpdateProfile
	}

	return nil
}

func (ur *UserRepository) GetPasswordByUserID(ctx context.Context, db DBTX, id int) (string, error) {
	query := `SELECT password FROM users WHERE id = $1`

	var password string
	err := db.QueryRow(ctx, query, id).Scan(&password)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return "", apperror.ErrUserNotFound
		}
		return "", apperror.ErrGetPassword
	}

	return password, nil
}

func (ur *UserRepository) UpdatePassword(ctx context.Context, db DBTX, id int, password string) error {
	query := `
		UPDATE users
		SET password = $1
		WHERE id = $2
	`

	_, err := db.Exec(ctx, query, password, id)
	if err != nil {
		log.Println(err.Error())
		return apperror.ErrUpdatePassword
	}

	return nil
}

func (ur *UserRepository) GetProfile(ctx context.Context, db DBTX, id int) (model.User, error) {
	query := `
		SELECT
		    id,
		    fullname,
		    email,
		    photo,
		    phone,
		    address,
		    created_at
		FROM users
		WHERE id = $1
	`

	row := db.QueryRow(ctx, query, id)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.Photo,
		&user.Phone,
		&user.Address,
		&user.CreatedAt,
	)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, apperror.ErrProfileNotFound
		}
		return model.User{}, apperror.ErrGetProfile
	}

	return user, nil
}

func (ur *UserRepository) InsertUser(ctx context.Context, db DBTX, req dto.InsertUserRequest, path string) error {
	query := `
		INSERT INTO
		    users (fullname, email, password, photo, phone, address, role)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.Exec(ctx, query, req.Fullname, req.Email, req.Password, path, req.Phone, req.Address, req.Role)
	if err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "duplicate") {
			return apperror.ErrEmailAlreadyExists
		}
		return apperror.ErrInsertUser
	}

	return nil
}
