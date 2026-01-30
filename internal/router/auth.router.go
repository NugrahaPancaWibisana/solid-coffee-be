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

func AuthRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := app.Group("/auth")

	authRepository := repository.NewAuthRepository()
	authService := service.NewAuthService(authRepository, rdb, db)
	authController := controller.NewAuthController(authService)

	authRouter.POST("/", authController.Login)
	authRouter.POST("/new", authController.Register)
	authRouter.DELETE("/", middleware.AuthMiddleware(), middleware.RBACMiddleware("user"), authController.Logout)
}
