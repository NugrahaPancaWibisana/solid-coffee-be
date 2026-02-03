package repository

import (
	"context"
	"log"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderRepo interface {
}
type OrderRepository struct {
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (o OrderRepository) GetPriceByMenuId(ctx context.Context, db DBTX, menuId int) (dto.MenuPriceResponse, error) {
	var price, discount float64
	var stock int

	sqlStr :=
		` SELECT 
				m.id, 
				p.price, 
				m.discount,
				m.stock
			FROM menus m
			JOIN products p ON p.id = m.product_id
			WHERE m.id = $1
		`

	values := []any{menuId}

	row := db.QueryRow(ctx, sqlStr, values...)

	if err := row.Scan(&menuId, &price, &discount, &stock); err != nil {
		log.Println(err.Error())
		return dto.MenuPriceResponse{}, err
	}

	return dto.MenuPriceResponse{
		Menu_Id:  menuId,
		Price:    price,
		Discount: discount,
		Stock:    stock,
	}, nil
}

func (o OrderRepository) CreateOrder(ctx context.Context, db DBTX, post dto.CreateOrder, userID int) (dto.CreateOrderResponse, error) {
	var orderId string
	var tax, total float64

	sqlStr := "INSERT INTO orders(shipping, tax, total, user_id, payment_id) VALUES (($1), ($2), ($3), ($4), ($5)) RETURNING id, tax, total"

	values := []any{post.Shipping, 0, 0, userID, post.Payment_Id}

	row := db.QueryRow(ctx, sqlStr, values...)

	if err := row.Scan(&orderId, &tax, &total); err != nil {
		log.Println(err.Error())
		return dto.CreateOrderResponse{}, err
	}

	return dto.CreateOrderResponse{
		Id_Order: orderId,
		Tax:      tax,
		Total:    total,
	}, nil
}

func (o OrderRepository) CreateDetailOrder(ctx context.Context, db DBTX, dt dto.CreateDetailOrder) (dto.CreateDetailOrderResponse, error) {

	var qty, menuId int
	var subtotal float64

	sqlStr := "INSERT INTO dt_order(order_id, qty, subtotal, menu_id, product_size_id, product_type_id) VALUES (($1), ($2), ($3), ($4), ($5), ($6)) RETURNING qty, subtotal, menu_id"

	values := []any{dt.OrderId, dt.Qty, dt.Subtotal, dt.MenuId, dt.ProductSizeId, dt.ProductTypeId}

	row := db.QueryRow(ctx, sqlStr, values...)

	if err := row.Scan(&qty, &subtotal, &menuId); err != nil {
		log.Println(err.Error())
		return dto.CreateDetailOrderResponse{}, err
	}

	return dto.CreateDetailOrderResponse{
		Qty:      qty,
		Subtotal: subtotal,
		MenuId:   menuId,
	}, nil
}

func (o OrderRepository) UpdateStockByIdMenu(ctx context.Context, db DBTX, updt dto.UpdateStock) (pgconn.CommandTag, error) {
	sqlStr := `
			UPDATE menus SET stock = $1 WHERE id = $2
			`
	values := []any{updt.Stock, updt.MenuId}
	return db.Exec(ctx, sqlStr, values...)
}

func (o OrderRepository) UpdateOrderById(ctx context.Context, db DBTX, updt dto.UpdateOrder) (pgconn.CommandTag, error) {
	sqlStr := `
			UPDATE orders SET tax = $1, total = $2 WHERE id = $3
			`
	values := []any{updt.Tax, updt.Total, updt.OrderId}
	return db.Exec(ctx, sqlStr, values...)
}

func (o OrderRepository) UpdateStatusByOrderId(ctx context.Context, db DBTX, updt dto.UpdateStatusOrder) (pgconn.CommandTag, error) {
	sqlStr := `
			UPDATE orders SET status = $1 WHERE id = $2
			`
	values := []any{strings.ToLower(updt.Status), updt.OrderId}
	return db.Exec(ctx, sqlStr, values...)
}

func (or *OrderRepository) AddReview(ctx context.Context, db DBTX, req dto.AddReview) error {
	query := "INSERT INTO reviews (rating, dt_orderid) VALUES ($1, $2)"

	_, err := db.Exec(ctx, query, req.Rating, req.DtOrderId)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
