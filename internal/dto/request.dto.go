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

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=8" example:"example123"`
	NewPassword string `json:"new_password" binding:"required,min=8,nefield=OldPassword" example:"example321"`
}

type PostProductsRequest struct {
	ProductName string  `form:"product_name,omitempty" json:"product_name"`
	Price       float32 `form:"price,omitempty" json:"price"`
	Description string  `form:"description,omitempty" json:"description"`
}

type PostImagesRequest struct {
	ImagesFile  []*multipart.FileHeader `form:"images_file,omitempty" json:"images_file"`
	Images_Name []string                `form:"images_name,omitempty" json:"images_name"`
}

type InsertUserRequest struct {
	Photo    *multipart.FileHeader `form:"photo"`
	Fullname string                `form:"fullname" binding:"required,min=3" example:"John Doe"`
	Email    string                `form:"email" binding:"required,email" example:"example123@gmail.com"`
	Phone    string                `form:"phone" binding:"required,min=3" example:"08123456789"`
	Password string                `form:"password" binding:"required,min=8" example:"example123"`
	Address  string                `form:"address" binding:"required,min=3" example:"Jakarta"`
	Role     string                `form:"role" binding:"required,oneof=user admin" example:"user"`
}

type UpdateProductsRequest struct {
	ProductName string  `form:"product_name,omitempty" json:"product_name"`
	Price       float32 `form:"price,omitempty" json:"price"`
	Description string  `form:"description,omitempty" json:"description"`
}

type UserParams struct {
	ID int `uri:"id"`
}

type UserQueries struct {
	Page string `form:"page"`
}

type ProductQueries struct {
	ID       string   `form:"id"`
	Category []string `form:"category"`
	Page     string   `form:"page"`
	Title    string   `form:"title"`
	Min      string   `form:"min"`
	Max      string   `form:"max"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"example123@gmail.com"`
}

type UpdateForgotPasswordRequest struct {
	Otp             string `json:"otp_code" binding:"required,max=6" example:"123456"`
	NewPassword     string `json:"password" binding:"required,min=8" example:"example123"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword" example:"example123"`
}

type CreateOrder struct {
	Shipping   string `json:"shipping" binding:"required"`
	Payment_Id int    `json:"payment_id" binding:"required"`
	// Status   string            `json:"status" binding:"required"`
	Menus []CreateMenuOrder `json:"menus" binding:"required"`
}

type CreateDetailOrder struct {
	OrderId       string  `json:"order_id"`
	MenuId        int     `json:"menu_id"`
	Qty           int     `json:"qty"`
	ProductSizeId int     `json:"product_size_id"`
	ProductTypeId int     `json:"product_type_id"`
	Subtotal      float64 `json:"subtotal"`
}

type CreateMenuOrder struct {
	MenuId        int `json:"menu_id"`
	Qty           int `json:"qty"`
	ProductSizeId int `json:"product_size_id"`
	ProductTypeId int `json:"product_type_id"`
}

type UpdateOrder struct {
	OrderId string  `json:"order_id" binding:"required"`
	Tax     float64 `json:"tax" binding:"required"`
	Total   float64 `json:"total" binding:"required"`
}

type UpdateStock struct {
	MenuId int `json:"menu_id"`
	Stock  int `json:"stock"`
}

type UpdateStatusOrder struct {
	OrderId string `json:"order_id" binding:"required"`
	Status  string `json:"status" binding:"required"`
}

type AddReview struct {
	DtOrderId int `json:"dt_orderid" binding:"required"`
	Rating    int `json:"rating" binding:"required"`
}

type OrderQueries struct {
	Page string `json:"page"`
}
