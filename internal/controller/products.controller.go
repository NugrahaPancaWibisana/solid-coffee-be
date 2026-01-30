package controller

import (
	"net/http"
	"strconv"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductsController struct {
	productService *service.ProductService
}

func NewProductsController(productService *service.ProductService) *ProductsController {
	return &ProductsController{
		productService: productService,
	}
}

// Get Products By Status godoc
// @Summary      Get All Products
// @Tags         products
// @Produce      json
// @Param        page		query int  true  "Pages"
// @Success      200  {object}  dto.Products
// @Failure 		 500 {object} dto.ResponseError
// @Router       /products/ [get]
func (p ProductsController) GetAllProducts(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("page"))
	data, err := p.productService.GetAllProducts(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Message: "Internal Server Error",
			Status:  "Internal Server Error",
			Error:   "internal server error",
		})
		return
	}

	totalPage, err := p.productService.GetTotalPage(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Message: "Internal Server Error",
			Status:  "Internal Server Error",
			Error:   "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Msg:     "OK",
		Success: true,
		Data:    []any{data},
		Meta: dto.PaginationMeta{
			Page:      id,
			TotalPage: totalPage,
		},
	})
}
