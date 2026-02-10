package dto

type ResponseSuccess struct {
	Status  string `json:"status" example:"Success"`
	Message string `json:"message" example:"Data retrieved successfully"`
}

type ResponseError struct {
	Status  string `json:"status" example:"Error"`
	Message string `json:"message" example:"Failed get data"`
	Error   string `json:"errors,omitempty" example:"failed get data"`
}

type LoginResponse struct {
	ResponseSuccess
	Data JWT `json:"data"`
}

type RegisterResponse struct {
	ResponseSuccess
}

type ProductResponse struct {
	ResponseSuccess
	Data []Products     `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page      int    `json:"page,omitempty"`
	TotalPage int    `json:"total_page,omitempty"`
	NextPage  string `json:"next_page,omitempty"`
	PrevPage  string `json:"prev_page,omitempty"`
}

type PostProductResponse struct {
	Id int `json:"id,omitempty"`
}

type UserProfileResponse struct {
	ResponseSuccess
	Data User
}

type UpdateProductResponse struct {
	Id int `json:"id,omitempty"`
}

type MenuPriceResponse struct {
	Menu_Id  int     `json:"menu_id,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Discount float64 `json:"discount,omitempty"`
	Stock    int     `json:"stock,omitempty"`
}
type CreateOrderResponse struct {
	Id_Order string  `json:"id,omitempty"`
	Tax      float64 `json:"tax,omitempty"`
	Total    float64 `json:"total,omitempty"`
}

type CreateDetailOrderResponse struct {
	Qty      int     `json:"qty,omitempty"`
	Subtotal float64 `json:"subtotal,omitempty"`
	MenuId   int     `json:"menu_id,omitempty"`
}

type DetailOrderResponse struct {
	Order_Id      string               `json:"order_id"`
	DateOrder     string               `json:"date_order"`
	FullName      string               `json:"fullname"`
	Address       string               `json:"address"`
	Phone         string               `json:"phone"`
	PaymentMethod string               `json:"payment_method"`
	Shipping      string               `json:"shipping"`
	Status        string               `json:"status"`
	Total         string               `json:"total"`
	DetailItem    []DetailItemResponse `json:"detail_item"`
}

type DetailItemResponse struct {
	Detail_Id   int      `json:"detail_id"`
	ItemName    string   `json:"item_name"`
	Qty         int      `json:"qty"`
	Images      []string `json:"image"`
	ProductSize string   `json:"product_size"`
	ProductType string   `json:"product_type"`
	Subtotal    string   `json:"subtotal"`
}
