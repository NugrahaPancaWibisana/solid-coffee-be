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
