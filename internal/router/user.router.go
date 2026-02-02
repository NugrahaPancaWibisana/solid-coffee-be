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

func UserRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	adminUserRouter := app.Group("/admin/user")
	userRouter := app.Group("/user")
	adminUserRouter.Use(middleware.AuthMiddleware(), middleware.RBACMiddleware("admin"))
	userRouter.Use(middleware.AuthMiddleware(), middleware.RBACMiddleware("user", "admin"))

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, rdb, db)
	userController := controller.NewUserController(userService)

	userRouter.GET("/", userController.GetProfile)
	userRouter.PATCH("/", userController.UpdateProfile)
	userRouter.PATCH("/password", userController.UpdatePassword)

	adminUserRouter.POST("/", userController.InsertUser)
	adminUserRouter.DELETE("/:id", userController.DeleteUser)
}
