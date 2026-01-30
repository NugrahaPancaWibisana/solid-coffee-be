package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")

		AllowOrigins := []string{"http://localhost:5173", "http://localhost:8080", "http://192.168.50.221:8080"}
		AllowHeaders := []string{"Origin", "Content-Type", "Authorization"}
		AllowMethods := []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions}

		if slices.Contains(AllowOrigins, origin) {
			ctx.Header("Access-Control-Allow-Origin", origin)
		}

		ctx.Header("Access-Control-Allow-Headers", strings.Join(AllowHeaders, ", "))
		ctx.Header("Access-Control-Allow-Methods", strings.Join(AllowMethods, ", "))

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
