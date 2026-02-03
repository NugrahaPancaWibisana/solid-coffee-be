package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/response"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/service"
	jwtutil "github.com/NugrahaPancaWibisana/solid-coffee-be/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type OrdersController struct {
	orderService *service.OrderService
}

func NewOrdersController(orderService *service.OrderService) *OrdersController {
	return &OrdersController{
		orderService: orderService,
	}
}

// Post Product godoc
// @Summary      Create order
// @Tags         orders
// @accept			 json
// @Produce      json
// @Param        movies	 body dto.CreateOrder  true  "Create order"
// @Success      200  {object}  dto.ResponseSuccess
// @Failure 		 500 {object} dto.ResponseError
// @Failure			 404 {object} dto.ResponseError
// @Failure			 400 {object} dto.ResponseError
// @Failure			 401 {object} dto.ResponseError
// @Router       /orders/ [post]
// @Security			BearerAuth
func (o OrdersController) CreateOrder(c *gin.Context) {

	var createOrder dto.CreateOrder

	token, isExist := c.Get("token")
	if !isExist {
		response.Error(c, http.StatusForbidden, "Forbidden Access")
		return
	}

	accessToken, _ := token.(jwtutil.JwtClaims)

	if err := c.ShouldBindJSON(&createOrder); err != nil {
		log.Println(err.Error())
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	userId := accessToken.UserID

	data, err := o.orderService.CreateOrder(c.Request.Context(), createOrder, userId)

	if err != nil {
		str := err.Error()
		if strings.Contains(str, "empty") {
			response.Error(c, http.StatusBadRequest, "Invalid Body")
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid Body")
		return
	}
	response.Success(c, http.StatusOK, "Order Created Successfully", dto.CreateOrderResponse{
		Id_Order: data.Id_Order,
	})
}

// UpdateProduct godoc
// @Summary      Update status
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        orders	 body dto.UpdateStatusOrder  true  "Update status order"
// @Success      200  {object}  dto.ResponseSuccess
// @Failure			401 {object} dto.ResponseError
// @Failure 		 400 {object} dto.ResponseError
// @Failure			 404 {object} dto.ResponseError
// @Failure 		 500 {object} dto.ResponseError
// @Router       /admin/orders/ [patch]
// @security 		 BearerAuth
func (o OrdersController) UpdateStatusOrder(c *gin.Context) {
	var updtStatus dto.UpdateStatusOrder

	_, isExist := c.Get("token")
	if !isExist {
		response.Error(c, http.StatusForbidden, "Forbidden Access")
		return
	}

	if err := c.ShouldBindJSON(&updtStatus); err != nil {
		log.Println(err.Error())
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := o.orderService.UpdateStatusByOrderId(c.Request.Context(), updtStatus); err != nil {
		str := err.Error()
		if strings.Contains(str, "empty") {
			response.Error(c, http.StatusBadRequest, "Invalid Body")
			return
		}
		if err.Error() == "no rows in result set" {
			response.Error(c, http.StatusNotFound, "Data Not Found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	response.Success(c, http.StatusOK, "Status Order Updated Successfully", nil)
}

// AddReview godoc
//
//	@Summary		Add review to order
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.AddReview	true	"Add review"
//	@Success		200		{object}	dto.ResponseSuccess
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		401		{object}	dto.ResponseError
//	@Failure		403		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/orders/review/ [post]
//	@Security		BearerAuth
func (o OrdersController) AddReview(ctx *gin.Context) {
	var req dto.AddReview

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "DtOrderId") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Detail Order ID is required")
			return
		}

		if strings.Contains(errStr, "Rating") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Rating is required")
			return
		}

		response.Error(ctx, http.StatusBadRequest, "Invalid request body")
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

	if err := o.orderService.AddReview(ctx.Request.Context(), req, accessToken.UserID, token[1]); err != nil {
		if err.Error() == "no rows in result set" {
			response.Error(ctx, http.StatusNotFound, "Order Not Found")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Success(ctx, http.StatusOK, "Review Added Successfully", nil)
}
