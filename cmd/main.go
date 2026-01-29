package main

import (
	"log"
	"os"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/config"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title						Solid Coffee Backend
// @version						1.0
// @description					Solid Coffee Backend RESTful API
// @host						localhost:8080
// @BasePath					/
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description					Type "Bearer" followed by a space and JWT token.
func main() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("Failed to Load env")
			return
		}
	}

	db, err := config.InitDB()
	if err != nil {
		log.Println("Failed to Connect to Database")
		return
	}

	rdb := config.InitRds()
	defer rdb.Close()

	app := gin.Default()

	router.Init(app, db, rdb)

	app.Run()
}
