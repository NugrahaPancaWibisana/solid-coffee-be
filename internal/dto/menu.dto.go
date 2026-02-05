package dto

type Menu struct {
	ID          int     `json:"id"`
	Discount    float64 `json:"discount"`
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Stock       int     `json:"stock"`
}
