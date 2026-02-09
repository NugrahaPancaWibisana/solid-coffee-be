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
	"strconv"
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
//	@Tags			Users
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

// UpdateProfileAdmin godoc
//
//	@Summary		Update user profile
//	@Description	Update authenticated user's profile information
//	@Tags			Admin User Management
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id			path		int		false	"user id"
//	@Param			photo		formData	file	false	"Profile photo"
//	@Param			fullname	formData	string	false	"Full name (min 3 chars)"
//	@Param			phone		formData	string	false	"Phone number (min 3 chars)"
//	@Param			address		formData	string	false	"Address (min 3 chars)"
//	@Success		200			{object}	dto.ResponseSuccess
//	@Failure		400			{object}	dto.ResponseError
//	@Failure		401			{object}	dto.ResponseError
//	@Router			/admin/user/{id} [patch]
//	@Security		BearerAuth
func (uc *UserController) UpdateProfileAdmin(ctx *gin.Context) {
	var param dto.UserParams
	if err := ctx.ShouldBindUri(&param); err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

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

	reqChange := dto.UpdateProfileRequest{
		Photo:    req.Photo,
		Fullname: req.Fullname,
		Phone:    req.Phone,
		Address:  req.Address,
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

	oldPath, err := uc.userService.UpdateProfileAdmin(ctx, reqChange, imagePath, param.ID, accessToken.UserID, token[1])
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
//	@Tags			Users
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

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusOK, "Password updated successfully", nil)
}

// GetProfile godoc
//
//	@Summary		Get user profile
//	@Description	Get authenticated user's profile information
//	@Tags			Users
//	@Produce		json
//	@path
//	@Success		200	{object}	dto.UserProfileResponse
//	@Failure		401	{object}	dto.ResponseError
//	@Router			/user/ [get]
//	@Security		BearerAuth
func (uc *UserController) GetProfile(ctx *gin.Context) {
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
	data, err := uc.userService.GetProfile(ctx, accessToken.UserID, token[1])
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusOK, "Profile retrieved successfully", data)
}

// InsertUser godoc
//
//	@Summary		Insert user profile
//	@Description	Insert authenticated user's profile information
//	@Tags			Admin User Management
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			photo		formData	file	true	"Profile photo"
//	@Param			fullname	formData	string	true	"Full name (min 3 chars)"
//	@Param			email		formData	string	true	"Email"
//	@Param			password	formData	string	true	"Password (min 8 chars)"
//	@Param			phone		formData	string	true	"Phone number (min 3 chars)"
//	@Param			address		formData	string	true	"Address (min 3 chars)"
//	@Param			role		formData	string	true	"role (user or admin)"
//	@Success		201			{object}	dto.ResponseSuccess
//	@Failure		400			{object}	dto.ResponseError
//	@Failure		401			{object}	dto.ResponseError
//	@Router			/admin/user/ [post]
//	@Security		BearerAuth
func (uc *UserController) InsertUser(ctx *gin.Context) {
	var req dto.InsertUserRequest
	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); err != nil {
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

		if strings.Contains(errStr, "Phone") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Phone number field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Phone") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Phone number must be at least 3 characters")
			return
		}

		if strings.Contains(errStr, "Address") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Address field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Address") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Address must be at least 3 characters")
			return
		}

		if strings.Contains(errStr, "Role") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Role field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Role") && strings.Contains(errStr, "oneof") {
			response.Error(ctx, http.StatusBadRequest, "Invalid role, allowed values: admin, user")
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

	if err := uc.userService.InsertUser(ctx, req, accessToken.UserID, imagePath, token[1]); err != nil {
		if errors.Is(err, apperror.ErrEmailAlreadyExists) || errors.Is(err, apperror.ErrInvalidEmailFormat) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		log.Println(err)
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusCreated, "User created successfully", nil)
}

// DeleteUser godoc
//
//	@Summary		Delete user profile
//	@Description	Delete authenticated user's profile information
//	@Tags			Admin User Management
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id	path		int	false	"user id"
//	@Success		201	{object}	dto.ResponseSuccess
//	@Failure		401	{object}	dto.ResponseError
//	@Router			/admin/user/{id} [delete]
//	@Security		BearerAuth
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	var param dto.UserParams
	if err := ctx.ShouldBindUri(&param); err != nil {
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
	if err := uc.userService.DeleteUser(ctx, accessToken.UserID, param.ID, token[1]); err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusOK, "User created successfully", nil)
}

// GetUser godoc
//
//	@Summary		Get all user profile
//	@Description	Get authenticated user's profile information
//	@Tags			Admin User Management
//	@Produce		json
//	@Param			page	query		string	false	"Page number"
//	@Success		200		{object}	dto.UserProfileResponse
//	@Failure		401		{object}	dto.ResponseError
//	@Router			/admin/user/ [get]
//	@Security		BearerAuth
func (uc *UserController) GetUsers(ctx *gin.Context) {
	var req dto.UserQueries
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	page := 1
	if req.Page != "" {
		page, _ = strconv.Atoi(req.Page)
		if page < 1 {
			page = 1
		}
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
	data, totalPage, err := uc.userService.GetUsers(ctx, req, accessToken.UserID, token[1])
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	var nextPage string
	var prevPage string

	if page < totalPage {
		nextPage = fmt.Sprintf("/admin/users?page=%d", page+1)
	}
	if page > 1 {
		prevPage = fmt.Sprintf("/admin/users?page=%d", page-1)
	}

	response.SuccessWithMeta(ctx, http.StatusOK, "Users data Retrieved Successfully", data,
		dto.PaginationMeta{
			Page:      page,
			TotalPage: totalPage,
			NextPage:  nextPage,
			PrevPage:  prevPage,
		},
	)
}
