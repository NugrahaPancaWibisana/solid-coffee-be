package dto

import "mime/multipart"

type Products struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Images_Name string  `json:"image_products"`
	Price       float32 `json:"price"`
	Discount    float32 `json:"discount"`
	Rating      float32 `json:"rating_product"`
}

type PostProducts struct {
	ProductName string  `form:"product_name,omitempty" json:"product_name"`
	Price       float32 `form:"price,omitempty" json:"price"`
	Description string  `form:"description,omitempty" json:"description"`
}

type PostImages struct {
	ImagesFile  []*multipart.FileHeader `form:"images_file,omitempty" json:"images_file"`
	Images_Name []string                `form:"images_name,omitempty" json:"images_name"`
}
