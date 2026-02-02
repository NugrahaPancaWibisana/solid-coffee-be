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

func (p ProductRepository) GetAllProduct(ctx context.Context, db DBTX, page int) ([]model.Products, error) {
	sqlStr :=
		`WITH avg_rating AS (
  		SELECT AVG(r.rating) AS "rating_product",
  		d.menu_id AS "idmenu"
  	FROM reviews r
  	JOIN dt_order d ON d.id = r.id
  	JOIN menus m ON m.id = d.menu_id
  	GROUP BY d.menu_id
	)

	SELECT
			p.id,
    	p.name,
    	string_agg(pi.image, ',') AS "image products",
    	p.price,
    	m.discount,
    	ar."rating_product"
  	FROM menus m
  	JOIN avg_rating ar ON ar."idmenu"= m.id
  	JOIN products p ON p.id = m.product_id
  	JOIN product_images pi ON pi.product_id = m.product_id
  	GROUP BY p.id, m.id, ar."rating_product" LIMIT 6 OFFSET 
	`

	offset := (page * 6) - 6
	spt := sqlStr + strconv.Itoa(offset)
	sql := spt

	rows, err := db.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var products []model.Products
	for rows.Next() {
		var product model.Products
		err := rows.Scan(&product.Id, &product.Name, &product.Images_Name, &product.Price, &product.Discount,
			&product.Rating)

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (p ProductRepository) GetTotalPage(ctx context.Context, db DBTX) (int, error) {
	sqlStr :=
		`WITH avg_rating AS (
  		SELECT AVG(r.rating) AS "rating_product",
  		d.menu_id AS "idmenu"
  	FROM reviews r
  	JOIN dt_order d ON d.id = r.id
  	JOIN menus m ON m.id = d.menu_id
  	GROUP BY d.menu_id
	)

	SELECT
			p.id,
    	p.name,
    	string_agg(pi.image, ',') AS "image products",
    	p.price,
    	m.discount,
    	ar."rating_product"
  	FROM menus m
  	JOIN avg_rating ar ON ar."idmenu"= m.id
  	JOIN products p ON p.id = m.product_id
  	JOIN product_images pi ON pi.product_id = m.product_id
  	GROUP BY p.id, m.id, ar."rating_product"
	`

	rows, err := db.Query(ctx, sqlStr)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	var products []model.Products
	for rows.Next() {
		var product model.Products
		err := rows.Scan(&product.Id, &product.Name, &product.Images_Name, &product.Price, &product.Discount,
			&product.Rating)

		if err != nil {
			log.Println(err.Error())
			return 0, err
		}
		products = append(products, product)
	}

	count := len(products)
	totalPage := math.Ceil(float64(count) / float64(6))

	return int(totalPage), nil
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
