package controller

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/service"
	jwtutil "github.com/NugrahaPancaWibisana/solid-coffee-be/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	page, _ := strconv.Atoi(c.Query("page"))
	var nextPage string
	var prevPage string
	data, err := p.productService.GetAllProducts(c.Request.Context(), page)
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

	if page < totalPage {
		nextPage = fmt.Sprintf("/products?page=%d", page+1)
		prevPage = fmt.Sprintf("/products?page=%d", page-1)
	}

	c.JSON(http.StatusOK, dto.ProductResponse{
		ResponseSuccess: dto.ResponseSuccess{
			Message: "Products Retrieved Successfully",
			Status:  "Success",
		},
		Data: data,
		Meta: dto.PaginationMeta{
			Page:      page,
			TotalPage: totalPage,
			NextPage:  nextPage,
			PrevPage:  prevPage,
		},
	})
}

// Post Product godoc
// @Summary      Post Products
// @Tags         products
// @accept			 multipart/form-data
// @Produce      json
// @Param        images_file	formData []file true  "Product Images"
// @Param        product_name	formData string true  "Products Name"
// @Param        price	formData number true  "Price"
// @Param				 description formData string true "Description"
// @Success      200  {object}  dto.ResponseSuccess
// @Failure 		 500 {object} dto.ResponseError
// @Failure			 404 {object} dto.ResponseError
// @Failure			 400 {object} dto.ResponseError
// @Router       /products/ [post]
// @Security			BearerAuth
func (p ProductsController) PostProducts(c *gin.Context) {
	const maxSize = 2 * 1024 * 1024
	var postImages dto.PostImages

	token, isExist := c.Get("token")
	if !isExist {
		c.AbortWithStatusJSON(http.StatusForbidden, dto.ResponseError{
			Message: "Forbidden Access",
			Status:  "403 Forbidden",
			Error:   "Access Denied",
		})
		return
	}

	accessToken, _ := token.(jwtutil.JwtClaims)

	if err := c.ShouldBindWith(&postImages, binding.FormMultipart); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Message: "Internal Server Error",
			Status:  "Internal Server Error",
			Error:   "internal server error",
		})
		return
	}

	if postImages.ImagesFile != nil {
		for key := range len(postImages.ImagesFile) {
			extPoster := path.Ext(postImages.ImagesFile[key].Filename)
			re := regexp.MustCompile("^[.](jpg|png)$")
			if !re.Match([]byte(extPoster)) {
				c.JSON(http.StatusBadRequest, dto.ResponseError{
					Message: "File have to be jpg or png",
					Error:   "Bad Request",
					Status:  "Bad Request",
				})
				return
			}
			//validasi ukuran
			if postImages.ImagesFile[key].Size > maxSize {
				c.JSON(http.StatusBadRequest, dto.ResponseError{
					Message: "File maximum 2 MB",
					Error:   "Bad Request",
					Status:  "Bad Request",
				})
				return
			}

			filenamePoster := fmt.Sprintf("%d_product_%d%s", time.Now().UnixNano(), accessToken.UserID, extPoster)
			postImages.Images_Name = append(postImages.Images_Name, filenamePoster)

			if e := c.SaveUploadedFile(postImages.ImagesFile[key], filepath.Join("public", "products", filenamePoster)); e != nil {
				log.Printf("error %v", e)
				c.JSON(http.StatusInternalServerError, dto.ResponseError{
					Message: "Internal Server Error",
					Status:  "Internal Server Error",
					Error:   "internal server error",
				})
				return
			}
		}
	}

	var newProduct dto.PostProducts

	if err := c.ShouldBindWith(&newProduct, binding.FormMultipart); err != nil {
		str := err.Error()
		if strings.Contains(str, "Field") {
			c.JSON(http.StatusBadRequest, dto.ResponseError{
				Message: "Invalid Body",
				Status:  "Bad Request",
				Error:   "invalid body",
			})
			return
		}
		if strings.Contains(str, "Empty") {
			c.JSON(http.StatusBadRequest, dto.ResponseError{
				Message: "Invalid Body",
				Status:  "Bad Request",
				Error:   "invalid body",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Message: "Internal Server Error",
			Status:  "Internal Server Error",
			Error:   "internal server error",
		})
		return
	}

	_, err := p.productService.PostProduct(c.Request.Context(), newProduct, postImages)

	if err != nil {
		str := err.Error()
		if strings.Contains(str, "empty") {
			c.JSON(http.StatusBadRequest, dto.ResponseError{
				Message: "Invalid Body",
				Status:  "Invalid Body",
				Error:   str,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ResponseError{
			Message: "Internal Server Error",
			Status:  "Internal Server Error",
			Error:   "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, dto.ResponseSuccess{
		Message: "Product Inserted",
		Status:  "Product Inserted",
	})
}
