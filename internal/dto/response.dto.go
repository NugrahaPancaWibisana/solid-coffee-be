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
