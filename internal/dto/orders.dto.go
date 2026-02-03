package dto

type Order struct {
	Order_Id string  `json:"order_id"`
	Date     string  `json:"date"`
	Item     string  `json:"item"`
	Status   string  `json:"status"`
	Total    float32 `json:"total"`
}
