package router

import (
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/controller"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/repository"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func ProductRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	productsRouter := app.Group("/products")

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository, rdb)
	productController := controller.NewProductsController(productService)

	productsRouter.GET("", productController.GetAllProducts)
}
