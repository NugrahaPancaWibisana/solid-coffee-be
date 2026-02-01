package dto

type Products struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Images_Name string  `json:"image_products"`
	Price       float32 `json:"price"`
	Discount    float32 `json:"discount"`
	Rating      float32 `json:"rating_product"`
}
