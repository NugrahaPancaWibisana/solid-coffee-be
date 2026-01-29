package controller

import "github.com/NugrahaPancaWibisana/solid-coffee-be/internal/service"

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}
