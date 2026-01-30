package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/service"
	jwtutil "github.com/NugrahaPancaWibisana/solid-coffee-be/pkg/jwt"
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

// Register godoc
//
//	@Summary		Register new user
//	@Description	Create a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"User registration data"
//	@Success		201		{object}	dto.RegisterResponse
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/auth/new/ [post]
func (ac *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "Fullname") && strings.Contains(errStr, "required") {
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Error:   "Bad Request",
				Message: "Fullname field cannot be empty",
				Status:  "error",
			})
			return
		}

		if strings.Contains(errStr, "Fullname") && strings.Contains(errStr, "min") {
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Error:   "Bad Request",
				Message: "Fullname must be at least 3 characters",
				Status:  "error",
			})
			return
		}

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

		if strings.Contains(errStr, "ConfirmPassword") && strings.Contains(errStr, "eqfield") {
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Error:   "Bad Request",
				Message: "Your password and confirmation password do not match.",
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

	err := ac.authService.Register(ctx, req)

	if err != nil {
		if errors.Is(err, apperror.ErrEmailAlreadyExists) || errors.Is(err, apperror.ErrInvalidEmailFormat) || errors.Is(err, apperror.ErrRegisterUser) {
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Error:   "Bad Request",
				Message: err.Error(),
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

	ctx.JSON(http.StatusCreated, dto.RegisterResponse{
		ResponseSuccess: dto.ResponseSuccess{
			Status:  "success",
			Message: "Registration successful",
		},
	})
}

// Logout godoc
//
//	@Summary		User logout
//	@Description	Logout user and invalidate token
//	@Tags			auth
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	dto.ResponseSuccess
//	@Failure		401	{object}	dto.ResponseError
//	@Failure		500	{object}	dto.ResponseError
//	@Router			/auth/ [delete]
//	@Security		BearerAuth
func (ac *AuthController) Logout(ctx *gin.Context) {
	claims, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ResponseError{
			Status:  "error",
			Error:   "Unauthorized",
			Message: "Unauthorized. Access token is missing, invalid, audience is incorrect, or has expired.",
		})
		return
	}

	jwtClaims, ok := claims.(jwtutil.JwtClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
			Status:  "error",
			Message: "Internal Server Error",
			Error:   "internal server error",
		})
		return
	}

	err := ac.authService.Logout(ctx, jwtClaims.UserID)
	if err != nil {
		if errors.Is(err, apperror.ErrSessionExpired) || errors.Is(err, apperror.ErrInvalidSession) || errors.Is(err, apperror.ErrLogoutFailed) {
			ctx.JSON(http.StatusBadRequest, dto.ResponseError{
				Error:   "Bad Request",
				Message: err.Error(),
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

	ctx.JSON(http.StatusOK, dto.ResponseSuccess{
		Status:  "success",
		Message: "Logout successful",
	})
}
