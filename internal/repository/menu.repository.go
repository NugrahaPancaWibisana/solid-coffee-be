package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/model"
	"github.com/jackc/pgx/v5"
)

type MenuRepo interface {
	CreateMenu(ctx context.Context, db DBTX, req dto.MenuRequest) error
	GetMenu(ctx context.Context, db DBTX, id int) (model.Menu, error)
	GetMenus(ctx context.Context, db DBTX, req dto.MenuParams) ([]model.Menu, error)
	GetTotalPage(ctx context.Context, db DBTX, req dto.MenuParams) (int, error)
	UpdateMenu(ctx context.Context, db DBTX, req dto.UpdateMenuRequest, id int) error
	DeleteMenu(ctx context.Context, db DBTX, id int) error
}

type MenuRepository struct{}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{}
}

func (mr *MenuRepository) CreateMenu(ctx context.Context, db DBTX, req dto.MenuRequest) error {
	query := `
		INSERT INTO
		    menus (discount, stock, product_id)
		VALUES
		    ($1, $2, $3)
	`

	_, err := db.Exec(ctx, query, req.Discount, req.Stock, req.ProductID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (mr *MenuRepository) GetMenu(ctx context.Context, db DBTX, id int) (model.Menu, error) {
	query := `
		SELECT
			m.id,
			m.discount,
			m.product_id,
			p.name AS product_name,
			m.stock
		FROM
			menus m
		JOIN products p ON p.id = m.product_id
		WHERE m.id = $1 AND m.deleted_at IS NULL
	`

	row := db.QueryRow(ctx, query, id)

	var menu model.Menu
	err := row.Scan(
		&menu.ID,
		&menu.Discount,
		&menu.ProductID,
		&menu.ProductName,
		&menu.Stock,
	)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Menu{}, apperror.ErrMenuNotFound
		}
		return model.Menu{}, apperror.ErrGetMenu
	}

	return menu, nil
}

func (mr *MenuRepository) GetMenus(ctx context.Context, db DBTX, req dto.MenuParams) ([]model.Menu, error) {
	var sb strings.Builder
	args := []any{}

	sb.WriteString(`
		SELECT
			m.id,
			m.discount,
			m.product_id,
			p.name AS product_name,
			m.stock
		FROM
			menus m
		JOIN products p ON p.id = m.product_id
		WHERE m.deleted_at IS NULL
	`)

	if req.Search != "" {
		fmt.Fprintf(&sb, " AND p.name ILIKE $%d", len(args)+1)
		args = append(args, "%"+req.Search+"%")
	}

	limit := 5
	offset := 0
	if req.Page != "" {
		page, _ := strconv.Atoi(req.Page)
		if page > 0 {
			offset = (page - 1) * limit
		}
	}

	fmt.Fprintf(&sb, " LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := db.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []model.Menu
	for rows.Next() {
		var menu model.Menu

		err := rows.Scan(
			&menu.ID,
			&menu.Discount,
			&menu.ProductID,
			&menu.ProductName,
			&menu.Stock,
		)
		if err != nil {
			return nil, err
		}

		menus = append(menus, menu)
	}

	return menus, nil
}

func (mr *MenuRepository) GetTotalPage(ctx context.Context, db DBTX, req dto.MenuParams) (int, error) {
	var sb strings.Builder
	args := []any{}

	sb.WriteString(`
		SELECT
			COUNT(DISTINCT m.id)
		FROM
			menus m
		JOIN products p ON p.id = m.product_id
		WHERE m.deleted_at IS NULL
	`)

	if req.Search != "" {
		fmt.Fprintf(&sb, " AND p.name ILIKE $%d", len(args)+1)
		args = append(args, "%"+req.Search+"%")
	}

	var totalMenus int
	err := db.QueryRow(ctx, sb.String(), args...).Scan(&totalMenus)
	if err != nil {
		return 0, err
	}

	itemsPerPage := 5
	totalPage := int(math.Ceil(float64(totalMenus) / float64(itemsPerPage)))

	return totalPage, nil
}

func (mr *MenuRepository) UpdateMenu(ctx context.Context, db DBTX, req dto.UpdateMenuRequest, id int) error {
	var sb strings.Builder
	sb.WriteString("UPDATE menus SET ")
	args := []any{}

	if req.ProductID != 0 {
		if len(args) > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "product_id = $%d", len(args)+1)
		args = append(args, req.ProductID)
	}

	if req.Discount != 0 {
		if len(args) > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "discount = $%d", len(args)+1)
		args = append(args, req.Discount)
	}

	if req.Stock != 0 {
		if len(args) > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "stock = $%d", len(args)+1)
		args = append(args, req.Stock)
	}

	if len(args) == 0 {
		return apperror.ErrNoFieldsToUpdate
	}

	fmt.Fprintf(&sb, " WHERE id = $%d AND deleted_at IS NULL", len(args)+1)
	args = append(args, id)

	_, err := db.Exec(ctx, sb.String(), args...)
	if err != nil {
		log.Println(err.Error())
		return apperror.ErrUpdateMenu
	}

	return nil
}

func (mr *MenuRepository) DeleteMenu(ctx context.Context, db DBTX, id int) error {
	query := "UPDATE menus SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL"

	_, err := db.Exec(ctx, query, id)
	if err != nil {
		log.Println(err.Error())
		return apperror.ErrDeleteMenu
	}

	return nil
}
