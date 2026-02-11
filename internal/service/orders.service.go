package service

import (
	"context"
	"errors"
	"log"
	"slices"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/cache"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type OrderService struct {
	orderRepository *repository.OrderRepository
	redis           *redis.Client
	db              *pgxpool.Pool
}

func NewOrderService(orderRepository *repository.OrderRepository, db *pgxpool.Pool, rdb *redis.Client) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		redis:           rdb,
		db:              db,
	}
}

func (o OrderService) CreateOrder(ctx context.Context, order dto.CreateOrder, userID int) (dto.CreateOrderResponse, error) {
	tx, err := o.db.Begin(ctx)
	if err != nil {
		log.Println(err)
		return dto.CreateOrderResponse{}, err
	}

	dataOrder, err := o.orderRepository.CreateOrder(ctx, tx, order, userID)
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}
	defer tx.Rollback(ctx)

	var totalSub float64

	for i := range len(order.Menus) {

		dataMenu, err := o.orderRepository.GetPriceByMenuId(ctx, tx, order.Menus[i].MenuId)

		var dt dto.CreateDetailOrder
		dt.OrderId = dataOrder.Id_Order
		dt.MenuId = order.Menus[i].MenuId
		discount := dataMenu.Price * dataMenu.Discount

		dt.ProductSizeId = order.Menus[i].ProductSizeId
		dt.ProductTypeId = order.Menus[i].ProductTypeId

		priceSize, err := o.orderRepository.GetProductSize(ctx, tx, dt.ProductSizeId)
		priceType, err := o.orderRepository.GetProductType(ctx, tx, dt.ProductTypeId)

		dt.Subtotal = ((dataMenu.Price - discount) * float64(order.Menus[i].Qty)) + float64(priceSize.Price) + float64(priceType.Price)
		dt.Qty = order.Menus[i].Qty

		currentStock := dataMenu.Stock - order.Menus[i].Qty

		if currentStock < 0 {
			return dto.CreateOrderResponse{}, errors.New("Stock Insufficient, Order Can't be Done !!")
		}

		stockUpdt := dataMenu.Stock - order.Menus[i].Qty
		totalSub = totalSub + dt.Subtotal

		var updtStock dto.UpdateStock
		updtStock.MenuId = order.Menus[i].MenuId
		updtStock.Stock = stockUpdt

		cmdx, err := o.orderRepository.UpdateStockByIdMenu(ctx, tx, updtStock)
		if err != nil {
			return dto.CreateOrderResponse{}, err
		}
		if cmdx.RowsAffected() == 0 {
			return dto.CreateOrderResponse{}, errors.New("no data updated")
		}

		_, e := o.orderRepository.CreateDetailOrder(ctx, tx, dt)
		if e != nil {
			return dto.CreateOrderResponse{}, err
		}
	}

	tax := totalSub * 0.1
	total := totalSub + tax

	var updtOrder dto.UpdateOrder

	updtOrder.OrderId = dataOrder.Id_Order
	updtOrder.Tax = tax
	updtOrder.Total = total

	cmd, err := o.orderRepository.UpdateOrderById(ctx, tx, updtOrder)
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}
	if cmd.RowsAffected() == 0 {
		return dto.CreateOrderResponse{}, errors.New("no data updated")
	}

	if e := tx.Commit(ctx); e != nil {
		log.Println("failed to commit", e.Error())
		return dto.CreateOrderResponse{}, e
	}

	response := dto.CreateOrderResponse{
		Id_Order: updtOrder.OrderId,
	}

	return response, nil
}

func (o OrderService) UpdateStatusByOrderId(ctx context.Context, sts dto.UpdateStatusOrder) error {
	status := []string{"pending, done, canceled"}
	isAvailable := slices.Contains(status, sts.Status)

	if !isAvailable {
		return errors.New("status is not valid")
	}

	cmd, err := o.orderRepository.UpdateStatusByOrderId(ctx, o.db, sts)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("no data deleted")
	}
	return nil
}

func (os *OrderService) AddReview(ctx context.Context, req dto.AddReview, id int, token string) error {
	if err := cache.CheckToken(ctx, os.redis, id, token); err != nil {
		return err
	}

	if err := os.orderRepository.AddReview(ctx, os.db, req); err != nil {
		return err
	}

	return nil
}

func (o *OrderService) GetAllOrderByAdmin(ctx context.Context, orderId string, status string, page int) ([]dto.Order, int, error) {

	totalPage, err := o.orderRepository.GetOrderTotalPages(ctx, o.db)
	if err != nil {
		return nil, 0, err
	}

	data, err := o.orderRepository.GetAllOrderByAdmin(ctx, o.db, status, orderId, page)
	if err != nil {
		return []dto.Order{}, 0, err
	}

	var response []dto.Order
	for _, v := range data {
		response = append(response, dto.Order{
			Order_Id: v.Order_Id,
			Date:     v.Date,
			Item:     v.Item,
			Status:   v.Status,
			Total:    v.Total,
		})
	}
	return response, totalPage, nil
}

func (o *OrderService) GetHistoryByUser(ctx context.Context, page int, userId int) ([]dto.History, int, error) {

	totalPage, err := o.orderRepository.GetHistoryTotalPages(ctx, o.db, userId)
	if err != nil {
		return nil, 0, err
	}

	data, err := o.orderRepository.GetHistoryByUser(ctx, o.db, page, userId)
	if err != nil {
		return []dto.History{}, 0, err
	}

	var response []dto.History
	for _, v := range data {
		response = append(response, dto.History{
			Order_Id: v.Order_Id,
			Date:     v.Date,
			Status:   v.Status,
			Total:    v.Total,
		})
	}
	return response, totalPage, nil
}

func (o *OrderService) GetDetailHistoryById(ctx context.Context, idOrder string) (dto.DetailOrderResponse, error) {
	var response dto.DetailOrderResponse

	data, err := o.orderRepository.GetOrderHistoryById(ctx, o.db, idOrder)
	if err != nil {
		return dto.DetailOrderResponse{}, err
	}

	dataDt, err := o.orderRepository.GetDetailOrderHistoryById(ctx, o.db, idOrder)
	imgStr := &[]string{}

	var resp []dto.DetailItemResponse
	for _, v := range dataDt {
		for _, img := range v.Image {
			*imgStr = append(*imgStr, img)
		}

		resp = append(resp, dto.DetailItemResponse{
			ItemName:    v.ItemName,
			ProductSize: v.ProductSize,
			ProductType: v.ProductType,
			Qty:         v.Qty,
			Subtotal:    v.Subtotal,
			Images:      *imgStr,
		})
	}

	response = dto.DetailOrderResponse{
		Order_Id:      data.Order_Id,
		DateOrder:     data.DateOrder,
		FullName:      data.FullName,
		Address:       data.Address,
		Phone:         data.Phone,
		PaymentMethod: data.PaymentMethod,
		Shipping:      data.Shipping,
		Status:        data.Status,
		Total:         data.Total,
		DetailItem:    resp,
	}

	return response, nil
}
