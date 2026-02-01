package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/response"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/service"
	jwtutil "github.com/NugrahaPancaWibisana/solid-coffee-be/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// UpdateProfile godoc
//
//	@Summary		Update user profile
//	@Description	Update authenticated user's profile information
//	@Tags			user
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			photo		formData	file	false	"Profile photo"
//	@Param			fullname	formData	string	false	"Full name (min 3 chars)"
//	@Param			phone		formData	string	false	"Phone number (min 3 chars)"
//	@Param			address		formData	string	false	"Address (min 3 chars)"
//	@Success		200			{object}	dto.ResponseSuccess
//	@Failure		400			{object}	dto.ResponseError
//	@Failure		401			{object}	dto.ResponseError
//	@Router			/user/ [patch]
//	@Security		BearerAuth
func (uc *UserController) UpdateProfile(ctx *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "no multipart boundary param in Content-Type") {
			response.Error(ctx, http.StatusBadRequest, "No fields to update")
			return
		}

		if strings.Contains(errStr, "Fullname") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Fullname must be at least 3 characters")
			return
		}

		if strings.Contains(errStr, "Phone") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Phone number must be at least 3 characters")
			return
		}

		if strings.Contains(errStr, "Address") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Address must be at least 3 characters")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	token := strings.Split(ctx.GetHeader("Authorization"), " ")
	if len(token) != 2 {
		response.Error(ctx, http.StatusUnauthorized, "Invalid Token")
		return
	}
	if token[0] != "Bearer" {
		response.Error(ctx, http.StatusUnauthorized, "Invalid Token")
		return
	}

	tokenData, _ := ctx.Get("token")
	accessToken, _ := tokenData.(jwtutil.JwtClaims)

	var imagePath string

	if req.Photo != nil {
		ext := strings.ToLower(path.Ext(req.Photo.Filename))
		re := regexp.MustCompile(`^\.(jpg|png)$`)
		if !re.MatchString(ext) {
			response.Error(ctx, http.StatusBadRequest, "File must be jpg or png")
			return
		}

		filename := fmt.Sprintf(
			"%d_profile_%d%s",
			time.Now().UnixNano(),
			accessToken.UserID,
			ext,
		)

		if err := ctx.SaveUploadedFile(
			req.Photo,
			filepath.Join("public", "profile", filename),
		); err != nil {
			log.Println(err.Error())
			response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		imagePath = fmt.Sprintf("/profile/%s", filename)
	}

	oldPath, err := uc.userService.UpdateProfile(ctx, req, imagePath, accessToken.UserID, token[1])
	if err != nil {
		if errors.Is(err, apperror.ErrNoFieldsToUpdate) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if oldPath != "" && imagePath != "" {
		oldFullPath := filepath.Join("public", oldPath)
		if err := os.Remove(oldFullPath); err != nil {
			log.Println("failed to delete old photo:", err.Error())
		}
	}

	response.Success(ctx, http.StatusOK, "Profile updated successfully", nil)
}

// UpdatePassword godoc
//
//	@Summary		Change user password
//	@Description	Update authenticated user's password
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UpdatePasswordRequest	true	"Edit password data"
//	@Success		200		{object}	dto.ResponseSuccess
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		401		{object}	dto.ResponseError
//	@Router			/user/password/ [patch]
//	@Security		BearerAuth
func (uc *UserController) UpdatePassword(ctx *gin.Context) {
	var req dto.UpdatePasswordRequest
	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "OldPassword") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Old password field cannot be empty")
			return
		}

		if strings.Contains(errStr, "OldPassword") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Old password must be at least 8 characters")
			return
		}

		if strings.Contains(errStr, "NewPassword") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "New password field cannot be empty")
			return
		}

		if strings.Contains(errStr, "NewPassword") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "New password must be at least 8 characters")
			return
		}

		if strings.Contains(errStr, "NewPassword") && strings.Contains(errStr, "nefield") {
			response.Error(ctx, http.StatusBadRequest, "The new password must be different from the current password")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	token := strings.Split(ctx.GetHeader("Authorization"), " ")
	if len(token) != 2 {
		response.Error(ctx, http.StatusUnauthorized, "Invalid Token")
		return
	}
	if token[0] != "Bearer" {
		response.Error(ctx, http.StatusUnauthorized, "Invalid Token")
		return
	}

	tokenData, _ := ctx.Get("token")
	accessToken, _ := tokenData.(jwtutil.JwtClaims)
	if err := uc.userService.UpdatePassword(ctx, req, accessToken.UserID, token[1]); err != nil {
		if errors.Is(err, apperror.ErrGetPassword) || errors.Is(err, apperror.ErrUpdatePassword) || errors.Is(err, apperror.ErrVerifyPassword) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}
		return
	}

	response.Success(ctx, http.StatusOK, "Password updated successfully", nil)
}
