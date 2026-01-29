package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	jwtutil "github.com/NugrahaPancaWibisana/solid-coffee-be/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := strings.Split(ctx.GetHeader("Authorization"), " ")
		if len(token) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
				Status:  "error",
				Message: "Unauthorized Access",
				Error:   "Invalid Token",
			})
			return
		}
		if token[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
				Status:  "error",
				Message: "Unauthorized Access",
				Error:   "Invalid Token",
			})
			return
		}

		var jc jwtutil.JwtClaims
		_, err := jc.VerifyToken(token[1])
		if err != nil {
			log.Println(err.Error())
			if errors.Is(err, jwt.ErrTokenExpired) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
					Status:  "error",
					Message: "Unauthorized Access",
					Error:   "Expired Token, Please Login Again",
				})
				return
			}
			if errors.Is(err, jwt.ErrTokenInvalidIssuer) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
					Status:  "error",
					Message: "Unauthorized Access",
					Error:   "Invalid Token, Please Login Again",
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ResponseError{
				Status:  "error",
				Message: "Internal Server Error",
				Error:   "internal server error",
			})
			return
		}
		ctx.Set("token", jc)
		ctx.Next()
	}
}
