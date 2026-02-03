package repository

import (
	"context"
	"log"
	"math"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/model"
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

func (o *OrderRepository) GetAllOrderByAdmin(ctx context.Context, db DBTX, page int) ([]model.Order, error) {

	sqlStr := `
		SELECT
			o.id,
			TO_CHAR(o.created_at, 'DD FMMonth YYYY') AS "date",
			STRING_AGG(CONCAT('• ' ,p.name, ' - ', dt.qty, 'x'), ', '),
			o.status,
			o.total
		FROM orders o
		JOIN dt_order dt ON dt.order_id = o.id
		JOIN menus m ON dt.menu_id = m.id
		JOIN products p ON p.id = m.product_id
		GROUP BY o.id LIMIT 5 OFFSET $1
	`

	offset := 0
	if page > 0 {
		offset = (page - 1) * 5
	}

	rows, err := db.Query(ctx, sqlStr, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var odr model.Order
		if err := rows.Scan(&odr.Order_Id, &odr.Date, &odr.Item, &odr.Status, &odr.Total); err != nil {
			return nil, err
		}
		orders = append(orders, odr)
	}

	return orders, rows.Err()

}

func (o *OrderRepository) GetOrderTotalPages(ctx context.Context, db DBTX) (int, error) {
	query := "SELECT COUNT(id) FROM orders"

	var order int
	err := db.QueryRow(ctx, query).Scan(&order)
	if err != nil {
		return 0, err
	}

	totalPage := int(math.Ceil(float64(order) / float64(5)))
	return totalPage, nil
}
