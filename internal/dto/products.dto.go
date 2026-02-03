package dto

type Products struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Images_Name string  `json:"image_products"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	Rating      float64 `json:"rating_product"`
}

type DetailProduct struct {
	IdProduct   int      `json:"id_product"`
	ProductName string   `json:"product_name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	IdImages    []string `json:"id_images"`
	Images      []string `json:"images"`
}

type DetailProductUser struct {
	IdProduct    int     `json:"id_product"`
	ProductName  string  `json:"product_name"`
	Images       string  `json:"images"`
	Price        float64 `json:"price"`
	Description  string  `json:"description"`
	Discount     float32 `json:"discount"`
	Rating       float64 `json:"rating"`
	Total_Review int     `json:"total_review"`
}
