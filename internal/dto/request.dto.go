package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"example123@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"example123"`
}

type RegisterRequest struct {
	Fullname        string `json:"fullname" binding:"required,min=3" example:"John Doe"`
	Email           string `json:"email" binding:"required,email" example:"example123@gmail.com"`
	Password        string `json:"password" binding:"required,min=8" example:"example123"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password" example:"example123"`
}
