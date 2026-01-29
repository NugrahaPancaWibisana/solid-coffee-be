package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user with email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	dto.LoginResponse
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		401		{object}	dto.ResponseError
//	@Router			/auth/ [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "required") {
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Error:   "Bad Request",
				Message: "Email field cannot be empty",
				Status:  "error",
			})
			return
		}

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "email") {
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Error:   "Bad Request",
				Message: "Email must be a valid email address",
				Status:  "error",
			})
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "required") {
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Error:   "Bad Request",
				Message: "Password field cannot be empty",
				Status:  "error",
			})
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "min") {
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Error:   "Bad Request",
				Message: "Password must be at least 8 characters",
				Status:  "error",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  "error",
			Message: "Internal Server Error",
			Error:   "internal server error",
		})
		return
	}

	data, err := ac.authService.Login(ctx, req)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) || errors.Is(err, apperror.ErrInvalidEmailFormat) {
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Error:   "Bad Request",
				Message: err.Error(),
				Status:  "error",
			})
			return
		}

		ctx.JSON(http.StatusUnauthorized, dto.ResponseError{
			Status:  "error",
			Message: "Invalid email or password",
			Error:   err.Error(),
		})
		return
	}

	token, err := ac.authService.GenerateJWT(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  "error",
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	ac.authService.WhitelistToken(ctx, data.ID, token)

	ctx.JSON(http.StatusOK, dto.LoginResponse{
		ResponseSuccess: dto.ResponseSuccess{
			Status:  "success",
			Message: "Login successful",
		},
		Data: dto.JWT{
			Token: token,
		},
	})
}
