package repository

import (
	"context"
	"log"
	"strconv"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (p ProductRepository) GetAllProduct(ctx context.Context, page int) ([]model.Products, error) {
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
	//spt := fmt.Sprintf(sqlStr, "%v", strconv.Itoa(offset))
	sql := spt

	rows, err := p.db.Query(ctx, sql)
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

func (p ProductRepository) GetTotalPage(ctx context.Context) (int, error) {
	var count int

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

	rows, err := p.db.Query(ctx, sqlStr)
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

	count = len(products)

	return count, nil
}
