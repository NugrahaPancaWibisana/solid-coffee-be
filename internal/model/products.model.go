package model

type Products struct {
	Id          int     `db:"id"`
	Name        string  `db:"name"`
	Images_Name string  `db:"image_products"`
	Price       float64 `db:"price"`
	Discount    float64 `db:"discount"`
	Rating      float64 `db:"rating_product"`
}

type DetailProduct struct {
	IdProduct   int      `db:"id_product"`
	ProductName string   `db:"product_name"`
	Description string   `db:"description"`
	Price       float64  `db:"price"`
	IdImages    []string `db:"id_images"`
	Images      []string `db:"images"`
}

type DetailProductUser struct {
	IdProduct    int     `db:"id_product"`
	ProductName  string  `db:"product_name"`
	Images       string  `db:"images"`
	Price        float64 `db:"price"`
	Description  string  `db:"description"`
	Discount     float32 `db:"discount"`
	Rating       float64 `db:"rating"`
	Total_Review int     `db:"total_review"`
}
