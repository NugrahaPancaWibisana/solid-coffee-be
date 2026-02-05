package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/response"
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
//	@Tags			Auth
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
			response.Error(ctx, http.StatusBadRequest, "Email field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "email") {
			response.Error(ctx, http.StatusBadRequest, "Email must be a valid email address")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Password field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	data, err := ac.authService.Login(ctx, req)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) || errors.Is(err, apperror.ErrInvalidEmailFormat) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(ctx, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := ac.authService.GenerateJWT(ctx, data)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	ac.authService.WhitelistToken(ctx, data.ID, token)

	response.Success(ctx, http.StatusOK, "Login successful", dto.JWT{Token: token})
}

// Register godoc
//
//	@Summary		Register new user
//	@Description	Create a new user account
//	@Tags			Auth
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
			response.Error(ctx, http.StatusBadRequest, "Fullname field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Fullname") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Fullname must be at least 3 characters")
			return
		}

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Email field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "email") {
			response.Error(ctx, http.StatusBadRequest, "Email must be a valid email address")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Password field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		if strings.Contains(errStr, "ConfirmPassword") && strings.Contains(errStr, "eqfield") {
			response.Error(ctx, http.StatusBadRequest, "Your password and confirmation password do not match.")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err := ac.authService.Register(ctx, req)

	if err != nil {
		if errors.Is(err, apperror.ErrEmailAlreadyExists) || errors.Is(err, apperror.ErrInvalidEmailFormat) || errors.Is(err, apperror.ErrRegisterUser) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusCreated, "Registration successful", nil)
}

// ForgotPassword godoc
//
//	@Summary		Request OTP for forgot password
//	@Description	Send OTP to registered email
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ForgotPasswordRequest	true	"Forgot password request"
//	@Success		200		{object}	dto.ResponseSuccess
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/auth/forgot-password/ [post]
func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest

	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Email field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "email") {
			response.Error(ctx, http.StatusBadRequest, "Email must be a valid email address")
			return
		}

		response.Error(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := ac.authService.ForgotPassword(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			response.Success(ctx, http.StatusOK, "If email is registered, OTP will be sent", nil)
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusOK, "If email is registered, OTP will be sent", nil)
}

// UpdateForgotPassword godoc
//
//	@Summary		Update password using OTP
//	@Description	Verify OTP and update user password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UpdateForgotPasswordRequest	true	"Update forgot password"
//	@Success		200		{object}	dto.ResponseSuccess
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/auth/forgot-password/update [post]
func (ac *AuthController) UpdateForgotPassword(ctx *gin.Context) {
	var req dto.UpdateForgotPasswordRequest

	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "Otp") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "OTP code is required")
			return
		}

		if strings.Contains(errStr, "NewPassword") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Password is required")
			return
		}

		if strings.Contains(errStr, "NewPassword") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		if strings.Contains(errStr, "ConfirmPassword") && strings.Contains(errStr, "eqfield") {
			response.Error(ctx, http.StatusBadRequest, "Passwords do not match")
			return
		}

		response.Error(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := ac.authService.UpdatePassword(ctx, req)
	if err != nil {
		if errors.Is(err, apperror.ErrOTPNotFound) || errors.Is(err, apperror.ErrOTPExpired) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusOK, "Password updated successfully", nil)
}
