package middleware

import (
	"net/http"
	"slices"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	jwtutil "github.com/NugrahaPancaWibisana/solid-coffee-be/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func RBACMiddleware(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, isExist := ctx.Get("token")
		if !isExist {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dto.ResponseError{
				Status:  "error",
				Message: "Forbidden Access",
				Error:   "Access Denied",
			})
			return
		}

		accessToken, ok := token.(jwtutil.JwtClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ResponseError{
				Status:  "error",
				Message: "Internal Server Error",
				Error:   "internal server error",
			})
			return
		}

		isAuthorized := slices.Contains(roles, accessToken.Role)
		if !isAuthorized {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dto.ResponseError{
				Status:  "error",
				Message: "Forbidden Access",
				Error:   "Access Denied",
			})
			return
		}

		ctx.Next()
	}
}
