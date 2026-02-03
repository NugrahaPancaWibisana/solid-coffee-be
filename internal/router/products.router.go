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

func ProductRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	adminProductsRouter := app.Group("/admin")

	productsRouter := app.Group("/products")
	productRepository := repository.NewProductRepository()
	productService := service.NewProductService(productRepository, db, rdb)
	productController := controller.NewProductsController(productService)

	productsRouter.GET("", productController.GetAllProducts)
	productsRouter.GET("/product-sizes", productController.GetAllProductSize)
	productsRouter.GET("/product-types", productController.GetAllProductType)

	adminProductsRouter.Use(middleware.AuthMiddleware())
	adminProductsRouter.GET("/products/:id", middleware.RBACMiddleware("admin"), productController.GetDetailProductById)
	adminProductsRouter.POST("/products", middleware.RBACMiddleware("admin"), productController.PostProducts)
	adminProductsRouter.PATCH("/products/:id", middleware.RBACMiddleware("admin"), productController.UpdateProduct)
	adminProductsRouter.DELETE("/products/:id", middleware.RBACMiddleware("admin"), productController.DeleteProductById)
	adminProductsRouter.DELETE("/products/image/:id", middleware.RBACMiddleware("admin"), productController.DeleteProductImageById)

}
