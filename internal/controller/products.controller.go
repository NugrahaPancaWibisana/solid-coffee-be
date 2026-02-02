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
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/response"
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
// @Summary      Get all products
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
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	totalPage, err := p.productService.GetTotalPage(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if page < totalPage {
		nextPage = fmt.Sprintf("/products?page=%d", page+1)
		prevPage = fmt.Sprintf("/products?page=%d", page-1)
	}

	response.SuccessWithMeta(c, http.StatusOK, "Products Retrieved Successfully", data,
		dto.PaginationMeta{
			Page:      page,
			TotalPage: totalPage,
			NextPage:  nextPage,
			PrevPage:  prevPage,
		},
	)
}

// Post Product godoc
// @Summary      Post product
// @Tags         products
// @accept			 multipart/form-data
// @Produce      json
// @Param        images_file	formData []file true  "Product images"
// @Param        product_name	formData string true  "Products name"
// @Param        price	formData number true  "Price"
// @Param				 description formData string true "Description"
// @Success      200  {object}  dto.ResponseSuccess
// @Failure 		 500 {object} dto.ResponseError
// @Failure			 404 {object} dto.ResponseError
// @Failure			 400 {object} dto.ResponseError
// @Router       /admin/products/ [post]
// @Security			BearerAuth
func (p ProductsController) PostProducts(c *gin.Context) {
	const maxSize = 2 * 1024 * 1024
	var postImages dto.PostImagesRequest

	token, isExist := c.Get("token")
	if !isExist {
		response.Error(c, http.StatusForbidden, "Forbidden Access")
		return
	}

	accessToken, _ := token.(jwtutil.JwtClaims)

	if err := c.ShouldBindWith(&postImages, binding.FormMultipart); err != nil {
		log.Println(err.Error())
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if postImages.ImagesFile != nil {
		for key := range len(postImages.ImagesFile) {
			extPoster := path.Ext(postImages.ImagesFile[key].Filename)
			re := regexp.MustCompile("^[.](jpg|png)$")
			if !re.Match([]byte(extPoster)) {
				response.Error(c, http.StatusBadRequest, "File have to be jpg or png")
				return
			}
			//validasi ukuran
			if postImages.ImagesFile[key].Size > maxSize {
				response.Error(c, http.StatusBadRequest, "File maximum 2 MB")
				return
			}

			filenamePoster := fmt.Sprintf("%d_product_%d%s", time.Now().UnixNano(), accessToken.UserID, extPoster)
			postImages.Images_Name = append(postImages.Images_Name, filenamePoster)

			if e := c.SaveUploadedFile(postImages.ImagesFile[key], filepath.Join("public", "products", filenamePoster)); e != nil {
				log.Printf("error %v", e)
				response.Error(c, http.StatusInternalServerError, "Internal Server Error")
				return
			}
		}
	}

	var newProduct dto.PostProductsRequest

	if err := c.ShouldBindWith(&newProduct, binding.FormMultipart); err != nil {
		str := err.Error()
		if strings.Contains(str, "Field") {
			response.Error(c, http.StatusBadRequest, "Invalid Body")
			return
		}
		if strings.Contains(str, "Empty") {
			response.Error(c, http.StatusBadRequest, "Invalid Body")
			return
		}
		response.Error(c, http.StatusBadRequest, "Internal Server Error")
		return
	}

	_, err := p.productService.PostProduct(c.Request.Context(), newProduct, postImages)

	if err != nil {
		str := err.Error()
		if strings.Contains(str, "empty") {
			response.Error(c, http.StatusBadRequest, "Invalid Body")
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid Body")
		return
	}
	response.Success(c, http.StatusOK, "Product Inserted", nil)
}

// UpdateProduct godoc
// @Summary      Update product
// @Tags         products
// @Accept       multipart/form-data
// @Produce      json
// @Param        id		path int  true  "Product Id"
// @Param        images_file	formData []file false  "Product images"
// @Param        product_name	formData string false  "Products name"
// @Param        price	formData number false  "Price"
// @Param				 description formData string false "Description"
// @Success      200  {object}  dto.ResponseSuccess
// @Failure 		 400 {object} dto.ResponseError
// @Failure 		 500 {object} dto.ResponseError
// @Router       /admin/products/{id} [patch]
// @security 		 BearerAuth
func (p ProductsController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	strId, _ := strconv.Atoi(id)

	const maxSize = 2 * 1024 * 1024
	var updateImages dto.PostImagesRequest

	token, isExist := c.Get("token")
	if !isExist {
		response.Error(c, http.StatusForbidden, "Forbidden Access")
		return
	}

	accessToken, _ := token.(jwtutil.JwtClaims)

	if err := c.ShouldBindWith(&updateImages, binding.FormMultipart); err != nil {
		log.Println(err.Error())
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if updateImages.ImagesFile != nil {
		for key := range len(updateImages.ImagesFile) {
			extPoster := path.Ext(updateImages.ImagesFile[key].Filename)
			re := regexp.MustCompile("^[.](jpg|png)$")
			if !re.Match([]byte(extPoster)) {
				response.Error(c, http.StatusBadRequest, "File have to be jpg or png")
				return
			}
			//validasi ukuran
			if updateImages.ImagesFile[key].Size > maxSize {
				response.Error(c, http.StatusBadRequest, "File maximum 2 MB")
				return
			}

			filenamePoster := fmt.Sprintf("%d_product_%d%s", time.Now().UnixNano(), accessToken.UserID, extPoster)
			updateImages.Images_Name = append(updateImages.Images_Name, filenamePoster)

			if e := c.SaveUploadedFile(updateImages.ImagesFile[key], filepath.Join("public", "products", filenamePoster)); e != nil {
				log.Printf("error %v", e)
				response.Error(c, http.StatusInternalServerError, "Internal Server Error")
				return
			}
		}
	}

	var updateProduct dto.UpdateProductsRequest

	if err := c.ShouldBindWith(&updateProduct, binding.FormMultipart); err != nil {
		str := err.Error()
		if strings.Contains(str, "Field") {
			response.Error(c, http.StatusBadRequest, "Invalid Body")
			return
		}
		if strings.Contains(str, "Empty") {
			response.Error(c, http.StatusBadRequest, "Invalid Body")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := p.productService.UpdateProduct(c.Request.Context(), updateProduct, updateImages, strId); err != nil {
		str := err.Error()
		if strings.Contains(str, "empty") {
			response.Error(c, http.StatusBadRequest, "Invalid Body")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Success(c, http.StatusOK, "Product Updated Successfully", nil)
}

// UpdateProduct godoc
// @Summary      Delete product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id		path int  true  "Product id"
// @Success      200  {object}  dto.ResponseSuccess
// @Failure 		 400 {object} dto.ResponseError
// @Failure 		 500 {object} dto.ResponseError
// @Router       /admin/products/{id} [delete]
// @security 		 BearerAuth
func (p ProductsController) DeleteProductById(c *gin.Context) {
	id := c.Param("id")
	strId, _ := strconv.Atoi(id)

	if err := c.ShouldBindUri(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := p.productService.DeleteProductById(c.Request.Context(), strId); err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	} else {
		response.Success(c, http.StatusOK, "Product Deleted Successfully", nil)
	}
}

// UpdateProduct godoc
// @Summary      Delete product image
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id		path int  true  "Image id"
// @Success      200  {object}  dto.ResponseSuccess
// @Failure 		 400 {object} dto.ResponseError
// @Failure 		 500 {object} dto.ResponseError
// @Router       /admin/products/image/{id} [delete]
// @security 		 BearerAuth
func (p ProductsController) DeleteProductImageById(c *gin.Context) {
	id := c.Param("id")
	strId, _ := strconv.Atoi(id)

	if err := c.ShouldBindUri(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := p.productService.DeleteProductImageById(c.Request.Context(), strId); err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	} else {
		response.Success(c, http.StatusOK, "Product Deleted Successfully", nil)
	}
}

// Get detail products by Id godoc
// @Summary      Get detail products by id
// @Tags         products
// @Produce      json
// @Param        id		path int  true  "Id"
// @Success      200  {object}  dto.Products
// @Failure 		 500 {object} dto.ResponseError
// @Router       /admin/products/{id} [get]
// @security		BearerAuth
func (p ProductsController) GetDetailProductById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindUri(id); err != nil {
		str := err.Error()
		if strings.Contains(str, "Empty") {
			response.Error(c, http.StatusBadRequest, "Invalid Body")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	data, err := p.productService.GetProductById(c.Request.Context(), id)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Products Retrieved Successfully", data, nil)

}
