package model

type Products struct {
	Id          int     `db:"id"`
	Name        string  `db:"name"`
	Images_Name string  `db:"image_products"`
	Price       float32 `db:"price"`
	Discount    float32 `db:"discount"`
	Rating      float32 `db:"rating_product"`
}
