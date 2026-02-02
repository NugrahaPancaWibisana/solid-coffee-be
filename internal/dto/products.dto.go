package dto

type Products struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Images_Name string  `json:"image_products"`
	Price       float32 `json:"price"`
	Discount    float32 `json:"discount"`
	Rating      float32 `json:"rating_product"`
}

type DetailProduct struct {
	IdProduct   int      `json:"id_product"`
	ProductName string   `json:"product_name"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	IdImages    []string `json:"id_images"`
	Images      []string `json:"images"`
}
