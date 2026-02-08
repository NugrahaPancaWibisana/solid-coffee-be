package dto

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID int    `json:"id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWT struct {
	Token string `json:"token"`
}
