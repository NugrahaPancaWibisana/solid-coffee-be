package dto

type Order struct {
	Order_Id string  `json:"order_id"`
	Date     string  `json:"date"`
	Item     string  `json:"item"`
	Status   string  `json:"status"`
	Total    float32 `json:"total"`
}

type History struct {
	Order_Id string  `json:"order_id"`
	Date     string  `json:"date"`
	Total    float32 `json:"total"`
	Status   string  `json:"status"`
}

type DetailOrder struct {
	Order_Id      string       `json:"order_id"`
	FullName      string       `json:"fullname"`
	Address       string       `json:"address"`
	Phone         string       `json:"phone"`
	PaymentMethod string       `json:"payment_method"`
	Shipping      string       `json:"shipping"`
	Status        string       `json:"status"`
	Total         string       `json:"total"`
	DetailItem    []DetailItem `json:"detail_item"`
}

type DetailItem struct {
	ItemName    string `json:"item_name"`
	Qty         int    `json:"qty"`
	ProductSize string `json:"product_size"`
	ProductType string `json:"product_type"`
	Price       string `json:"price"`
}
