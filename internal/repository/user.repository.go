package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/jackc/pgx/v5"
)

type UserRepo interface {
	GetPhoto(ctx context.Context, db DBTX, id int) (string, error)
	UpdateProfile(ctx context.Context, db DBTX, req dto.UpdateProfileRequest, path string, id int) error
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
