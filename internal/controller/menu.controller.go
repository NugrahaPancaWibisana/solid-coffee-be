package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/response"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/service"
	jwtutil "github.com/NugrahaPancaWibisana/solid-coffee-be/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MenuController struct {
	menuService *service.MenuService
}

func NewMenuController(menuService *service.MenuService) *MenuController {
	return &MenuController{menuService: menuService}
}

// CreateMenu godoc
//
//	@Summary		Create menu
//	@Description	Create a new menu item
//	@Tags			Admin Menu Management
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.MenuRequest	true	"Menu data"
//	@Success		201		{object}	dto.ResponseSuccess
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		401		{object}	dto.ResponseError
//	@Router			/admin/menu/ [post]
//	@Security		BearerAuth
func (mc *MenuController) CreateMenu(ctx *gin.Context) {
	var req dto.MenuRequest
	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "ProductID") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Product ID field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Stock") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Stock field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Discount") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Discount must be at least 0")
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

	if err := mc.menuService.CreateMenu(ctx, req, accessToken.UserID, token[1]); err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusCreated, "Menu created successfully", nil)
}

// GetMenu godoc
//
//	@Summary		Get menu by ID
//	@Description	Get menu item details by ID
//	@Tags			Admin Menu Management
//	@Produce		json
//	@Param			id	path		int	true	"Menu ID"
//	@Success		200	{object}	dto.ResponseSuccess
//	@Failure		401	{object}	dto.ResponseError
//	@Failure		404	{object}	dto.ResponseError
//	@Router			/admin/menu/{id} [get]
//	@Security		BearerAuth
func (mc *MenuController) GetMenu(ctx *gin.Context) {
	var param dto.MenuURIParam
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

	data, err := mc.menuService.GetMenu(ctx, accessToken.UserID, param.ID, token[1])
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusOK, "Menu retrieved successfully", data)
}

// GetMenus godoc
//
//	@Summary		Get all menus
//	@Description	Get all menu items with pagination
//	@Tags			Admin Menu Management
//	@Produce		json
//	@Param			page	query		string	false	"Page number"
//	@Success		200		{object}	dto.ResponseSuccess
//	@Failure		401		{object}	dto.ResponseError
//	@Router			/admin/menu/ [get]
//	@Security		BearerAuth
func (mc *MenuController) GetMenus(ctx *gin.Context) {
	var req dto.MenuParams
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

	data, totalPage, err := mc.menuService.GetMenus(ctx, req, accessToken.UserID, 0, token[1])
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	var nextPage string
	var prevPage string

	if page < totalPage {
		nextPage = fmt.Sprintf("/menu?page=%d", page+1)
	}
	if page > 1 {
		prevPage = fmt.Sprintf("/menu?page=%d", page-1)
	}

	response.SuccessWithMeta(ctx, http.StatusOK, "Menus retrieved successfully", data,
		dto.PaginationMeta{
			Page:      page,
			TotalPage: totalPage,
			NextPage:  nextPage,
			PrevPage:  prevPage,
		},
	)
}

// UpdateMenu godoc
//
//	@Summary		Update menu
//	@Description	Update menu item by ID
//	@Tags			Admin Menu Management
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Menu ID"
//	@Param			request	body		dto.UpdateMenuRequest	true	"Menu data"
//	@Success		200		{object}	dto.ResponseSuccess
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		401		{object}	dto.ResponseError
//	@Router			/admin/menu/{id} [patch]
//	@Security		BearerAuth
func (mc *MenuController) UpdateMenu(ctx *gin.Context) {
	var param dto.MenuURIParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	var req dto.UpdateMenuRequest
	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "Stock") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Stock must be at least 0")
			return
		}

		if strings.Contains(errStr, "Discount") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Discount must be at least 0")
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

	if err := mc.menuService.UpdateMenu(ctx, req, accessToken.UserID, param.ID, token[1]); err != nil {
		if errors.Is(err, apperror.ErrNoFieldsToUpdate) {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusOK, "Menu updated successfully", nil)
}

// DeleteMenu godoc
//
//	@Summary		Delete menu
//	@Description	Delete menu item by ID
//	@Tags			Admin Menu Management
//	@Produce		json
//	@Param			id	path		int	true	"Menu ID"
//	@Success		200	{object}	dto.ResponseSuccess
//	@Failure		401	{object}	dto.ResponseError
//	@Router			/admin/menu/{id} [delete]
//	@Security		BearerAuth
func (mc *MenuController) DeleteMenu(ctx *gin.Context) {
	var param dto.MenuURIParam
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

	if err := mc.menuService.DeleteMenu(ctx, accessToken.UserID, param.ID, token[1]); err != nil {
		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(ctx, http.StatusOK, "Menu deleted successfully", nil)
}
