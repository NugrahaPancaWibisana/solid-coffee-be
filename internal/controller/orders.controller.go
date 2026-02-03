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
