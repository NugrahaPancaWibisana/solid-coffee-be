package router

import (
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/controller"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/middleware"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/repository"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func OrderRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	adminOrdersRouter := app.Group("/admin")
	adminOrdersRouter.Use(middleware.AuthMiddleware())

	ordersRouter := app.Group("/orders")
	ordersRepository := repository.NewOrderRepository()
	ordersService := service.NewOrderService(ordersRepository, db, rdb)
	ordersController := controller.NewOrdersController(ordersService)
	ordersRouter.Use(middleware.AuthMiddleware())

	ordersRouter.POST("/", middleware.RBACMiddleware("user"), ordersController.CreateOrder)
	adminOrdersRouter.PATCH("/orders/", middleware.RBACMiddleware("admin"), ordersController.UpdateStatusOrder)
}
