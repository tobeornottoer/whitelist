package main

import (
	"log"
	"os"
	"whitelist/logger"
	"whitelist/route"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var RootPath,_ = os.Getwd()

func main(){
	loadEnv()
	gin.DisableConsoleColor()
	gin.SetMode(os.Getenv("RUN_MODE"))
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(logger.Logger())
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	route.Set(router)
	router.Run(os.Getenv("LISTEN_PORT"))
}

func loadEnv(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}