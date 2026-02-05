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

func MenuRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	menuRouter := app.Group("/admin/menu")
	menuRouter.Use(middleware.AuthMiddleware(), middleware.RBACMiddleware("admin"))

	menuRepository := repository.NewMenuRepository()
	menuService := service.NewMenuService(menuRepository, rdb, db)
	menuController := controller.NewMenuController(menuService)

	menuRouter.GET("/", menuController.GetMenus)
	menuRouter.GET("/:id", menuController.GetMenu)
	menuRouter.POST("/", menuController.CreateMenu)
	menuRouter.PATCH("/:id", menuController.UpdateMenu)
	menuRouter.DELETE("/:id", menuController.DeleteMenu)
}