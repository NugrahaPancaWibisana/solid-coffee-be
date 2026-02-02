package repository

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProductRepo interface {
	GetAllProduct(ctx context.Context, db DBTX, page int) ([]model.Products, error)
	GetTotalPage(ctx context.Context, db DBTX) (int, error)
	PostProduct(ctx context.Context, db DBTX, post dto.PostProductsRequest) (dto.PostProductResponse, error)
}
type ProductRepository struct {
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (pr *ProductRepository) GetProducts(ctx context.Context, db DBTX, req dto.ProductQueries) ([]model.Products, error) {
	var sb strings.Builder
	args := []any{}
	argCount := 1

	sb.WriteString(`
		WITH avg_rating AS (
			SELECT 
				AVG(r.rating) AS rating_product,
				d.menu_id AS idmenu
			FROM reviews r
			JOIN dt_order d ON d.id = r.dt_orderid
			GROUP BY d.menu_id
		)
		SELECT DISTINCT ON (p.id)
			p.id,
			p.name,
			STRING_AGG(DISTINCT pi.image, ',') AS image_products,
			p.price,
			m.discount,
			COALESCE(ar.rating_product, 0) AS rating_product
		FROM products p
		JOIN menus m ON m.product_id = p.id
		LEFT JOIN avg_rating ar ON ar.idmenu = m.id
		LEFT JOIN product_images pi ON pi.product_id = p.id
		LEFT JOIN product_categories pc ON pc.product_id = p.id
		LEFT JOIN categories c ON c.id = pc.category_id
		WHERE p.deleted_at IS NULL AND m.deleted_at IS NULL
	`)

	if len(req.Category) > 0 {
		placeholders := []string{}
		for _, cat := range req.Category {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argCount))
			args = append(args, cat)
			argCount++
		}
		fmt.Fprintf(&sb, " AND c.name IN (%s)", strings.Join(placeholders, ","))
	}

	if req.Title != "" {
		fmt.Fprintf(&sb, " AND p.name ILIKE $%d", argCount)
		args = append(args, "%"+req.Title+"%")
		argCount++
	}

	if req.Min != "" {
		fmt.Fprintf(&sb, " AND (p.price - (p.price * m.discount / 100)) >= $%d", argCount)
		args = append(args, req.Min)
		argCount++
	}

	if req.Max != "" {
		fmt.Fprintf(&sb, " AND (p.price - (p.price * m.discount / 100)) <= $%d", argCount)
		args = append(args, req.Max)
		argCount++
	}

	sb.WriteString(" GROUP BY p.id, p.name, p.price, m.discount, ar.rating_product")
	sb.WriteString(" ORDER BY p.id")

	limit := 6
	offset := 0
	if req.Page != "" {
		page, _ := strconv.Atoi(req.Page)
		if page > 0 {
			offset = (page - 1) * limit
		}
	}

	fmt.Fprintf(&sb, " LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	query := sb.String()

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Products
	for rows.Next() {
		var p model.Products

		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Images_Name,
			&p.Price,
			&p.Discount,
			&p.Rating,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, rows.Err()
}

func (pr *ProductRepository) GetTotalPage(ctx context.Context, db DBTX, req dto.ProductQueries) (int, error) {
	var sb strings.Builder
	args := []any{}
	argCount := 1

	sb.WriteString(`
		SELECT COUNT(DISTINCT p.id)
		FROM products p
		JOIN menus m ON m.product_id = p.id
		LEFT JOIN product_images pi ON pi.product_id = p.id
		LEFT JOIN product_categories pc ON pc.product_id = p.id
		LEFT JOIN categories c ON c.id = pc.category_id
		WHERE p.deleted_at IS NULL AND m.deleted_at IS NULL
	`)

	if len(req.Category) > 0 {
		placeholders := []string{}
		for _, cat := range req.Category {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argCount))
			args = append(args, cat)
			argCount++
		}
		fmt.Fprintf(&sb, " AND c.name IN (%s)", strings.Join(placeholders, ","))
	}

	if req.Title != "" {
		fmt.Fprintf(&sb, " AND p.name ILIKE $%d", argCount)
		args = append(args, "%"+req.Title+"%")
		argCount++
	}

	if req.Min != "" {
		fmt.Fprintf(&sb, " AND (p.price - (p.price * m.discount / 100)) >= $%d", argCount)
		args = append(args, req.Min)
		argCount++
	}

	if req.Max != "" {
		fmt.Fprintf(&sb, " AND (p.price - (p.price * m.discount / 100)) <= $%d", argCount)
		args = append(args, req.Max)
		argCount++
	}

	query := sb.String()

	var totalProducts int
	err := db.QueryRow(ctx, query, args...).Scan(&totalProducts)
	if err != nil {
		return 0, err
	}

	itemsPerPage := 6
	totalPage := int(math.Ceil(float64(totalProducts) / float64(itemsPerPage)))

	return totalPage, nil
}

func (p ProductRepository) PostProduct(ctx context.Context, db DBTX, post dto.PostProductsRequest) (dto.PostProductResponse, error) {
	var idProduct int
	sqlStr := "INSERT INTO products (name, price, description) VALUES (($1), ($2), ($3)) RETURNING id"
	values := []any{post.ProductName, post.Price, post.Description}

	row := db.QueryRow(ctx, sqlStr, values...)
	if err := row.Scan(&idProduct); err != nil {
		log.Println(err.Error())
		return dto.PostProductResponse{}, err
	}

	return dto.PostProductResponse{
		Id: idProduct,
	}, nil
}

func (p ProductRepository) PostImages(ctx context.Context, db DBTX, idProduct int, postImages string) (pgconn.CommandTag, error) {
	sqlStr := "INSERT INTO product_images (image, product_id) VALUES (($1), ($2))"
	values := []any{postImages, idProduct}
	return db.Exec(ctx, sqlStr, values...)
}

func (p ProductRepository) UpdateProduct(ctx context.Context, db DBTX, update dto.UpdateProductsRequest, id int) (pgconn.CommandTag, error) {
	var sql strings.Builder
	values := []any{}
	valuesAll := []any{}

	sql.WriteString("UPDATE products SET")
	if update.ProductName != "" {
		fmt.Fprintf(&sql, " name=$%d", len(values)+1)
		values = append(values, update.ProductName)
		valuesAll = append(valuesAll, &update.ProductName)
	}
	if update.Price != 0 {
		if len(values) > 0 {
			sql.WriteString(",")
		}
		fmt.Fprintf(&sql, " price=$%d", len(values)+1)
		values = append(values, update.Price)
		valuesAll = append(valuesAll, &update.Price)
	}
	if update.Description != "" {
		if len(values) > 0 {
			sql.WriteString(",")
		}
		fmt.Fprintf(&sql, " description=$%d", len(values)+1)
		values = append(values, update.Description)
		valuesAll = append(valuesAll, &update.Description)
	}
	if update.ProductName != "" || update.Price != 0 || update.Description != "" {
		sql.WriteString(" WHERE ")
		fmt.Fprintf(&sql, "id=%d", id)
	}

	sqlStr := sql.String()

	return db.Exec(ctx, sqlStr, valuesAll...)
}

func (p ProductRepository) DeleteProductById(ctx context.Context, db DBTX, idProduct int) (pgconn.CommandTag, error) {
	sqlStr := "UPDATE products SET deleted_at = NOW() WHERE id = $1"
	values := idProduct
	return db.Exec(ctx, sqlStr, values)
}

func (p ProductRepository) DeleteProductImage(ctx context.Context, db DBTX, idProduct int) (pgconn.CommandTag, error) {
	sqlStr := "UPDATE product_images SET deleted_at = NOW() WHERE product_id = $1"
	values := idProduct
	return db.Exec(ctx, sqlStr, values)
}

func (p ProductRepository) DeleteProductImageById(ctx context.Context, db DBTX, idImage int) (pgconn.CommandTag, error) {
	sqlStr := "UPDATE product_images SET deleted_at = NOW() WHERE id = $1"
	values := idImage
	return db.Exec(ctx, sqlStr, values)
}

func (p ProductRepository) GetProductById(ctx context.Context, db DBTX, idProduct int) (model.DetailProduct, error) {
	sqlStr := `
			SELECT p.id, 
			p.name, 
			p.description, 
			p.price, 
			array_agg(CAST(pi.id AS VARCHAR(50))), 
			array_agg(pi.image) 
			FROM products p 
			JOIN product_images pi ON pi.product_id = p.id 
			WHERE pi.product_id = $1 
			GROUP BY p.id`

	values := []any{idProduct}
	row := db.QueryRow(ctx, sqlStr, values...)

	var prdDetail model.DetailProduct

	if err := row.Scan(&prdDetail.IdProduct, &prdDetail.ProductName, &prdDetail.Description, &prdDetail.Price, &prdDetail.IdImages, &prdDetail.Images); err != nil {
		log.Println(err.Error())
		return model.DetailProduct{}, err
	}

	return prdDetail, nil
}
