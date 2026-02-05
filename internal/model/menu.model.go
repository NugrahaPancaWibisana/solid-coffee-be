package model

type Menu struct {
	ID          int     `db:"id"`
	Discount    float64 `db:"discount"`
	ProductID   int     `db:"product_id"`
	ProductName string  `db:"name"`
	Stock       int     `db:"stock"`
}
