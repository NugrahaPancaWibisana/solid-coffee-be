package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	_ "github.com/NugrahaPancaWibisana/solid-coffee-be/docs"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	app.Use(middleware.CORSMiddleware())
	AuthRouter(app, db, rdb)
	UserRouter(app, db, rdb)
	ProductRouter(app, db, rdb)
	OrderRouter(app, db, rdb)
	MenuRouter(app, db, rdb)

	app.Static("/static/img", "public")

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
