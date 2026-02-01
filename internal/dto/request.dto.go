package dto

import "mime/multipart"

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

type UpdateProfileRequest struct {
	Photo    *multipart.FileHeader `form:"photo"`
	Fullname string                `form:"fullname" binding:"omitempty,min=3" example:"John Doe"`
	Phone    string                `form:"phone" binding:"omitempty,min=3" example:"08123456789"`
	Address  string                `form:"address" binding:"omitempty,min=3" example:"Jakarta"`
}
