package response

import (
	"net/http"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, statusCode int, message string, data any) {
	res := dto.ResponseSuccess{
		Status:  "success",
		Message: message,
	}

	if data != nil {
		ctx.JSON(statusCode, struct {
			dto.ResponseSuccess
			Data any `json:"data"`
		}{
			ResponseSuccess: res,
			Data:            data,
		})
		return
	}

	ctx.JSON(statusCode, res)
}

func SuccessWithMeta(ctx *gin.Context, statusCode int, message string, data any, meta any) {
	ctx.JSON(statusCode, struct {
		dto.ResponseSuccess
		Data any `json:"data"`
		Meta any `json:"meta"`
	}{
		ResponseSuccess: dto.ResponseSuccess{
			Status:  "success",
			Message: message,
		},
		Data: data,
		Meta: meta,
	})
}

func Error(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, dto.ResponseError{
		Status:  "error",
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}
