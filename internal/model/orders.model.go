package model

type Order struct {
	Order_Id string  `db:"order_id"`
	Date     string  `db:"date"`
	Item     string  `db:"item"`
	Status   string  `db:"status"`
	Total    float32 `db:"total"`
}

type History struct {
	Order_Id string  `db:"order_id"`
	Date     string  `db:"date"`
	Total    float32 `db:"total"`
	Status   string  `db:"status"`
}

type DetailOrder struct {
	Order_Id      string       `db:"order_id"`
	FullName      string       `db:"fullname"`
	Address       string       `db:"address"`
	Phone         string       `db:"phone"`
	PaymentMethod string       `db:"payment_method"`
	Shipping      string       `db:"shipping"`
	Status        string       `db:"status"`
	Total         string       `db:"total"`
	DetailItem    []DetailItem `db:"detail_item"`
}

type DetailItem struct {
	ItemName    string `db:"item_name"`
	Qty         int    `db:"qty"`
	ProductSize string `db:"product_size"`
	ProductType string `db:"product_type"`
	Price       string `db:"price"`
}
