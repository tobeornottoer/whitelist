package main

import (
	"log"
	"net/http"
	"os"
	"whitelist/logger"
	"whitelist/route"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var RootPath,_ = os.Getwd()

func main(){
	loadEnv()
	gin.DisableConsoleColor()
	gin.SetMode(os.Getenv("RUN_MODE"))
	router := gin.Default()
	router.Use(Cors(),logger.Logger())
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusOK, "")
			c.Abort()
			return
		}
		c.Next()
	}
}