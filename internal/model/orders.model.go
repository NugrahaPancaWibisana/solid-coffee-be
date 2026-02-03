package model

type Order struct {
	Order_Id string  `db:"order_id"`
	Date     string  `db:"date"`
	Item     string  `db:"item"`
	Status   string  `db:"status"`
	Total    float32 `db:"total"`
}
