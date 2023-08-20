package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	logger = log.New()
	RuntimeLog *log.Entry
)

func init(){
	logger.SetFormatter(&log.JSONFormatter{})
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	logFile := os.Getenv("LOG_FILE")
	if logFile != "" {
		src, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err == nil {
			logger.Out = src
		}
	} else {
		logger.Out = os.Stdout
	}

	fmt.Print(logger.Out)
	logger.SetLevel(log.DebugLevel)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context){
		startTime 	:= time.Now()
		Runtime(c)
		c.Next()
		endTime 	:= time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod 	:= c.Request.Method
		reqUri 		:= c.Request.RequestURI
		statusCode 	:= c.Writer.Status()
		clientIP 	:= c.ClientIP()
		logger.WithFields(log.Fields{
			"request_method"	: reqMethod,
			"request_uri"		: reqUri,
			"status_code"		: statusCode,
			"client_ip"			: clientIP,
			"latency_time"		: latencyTime,
		}).Info()
	}
}

func Runtime(c *gin.Context) {
	RuntimeLog = logger.WithFields(log.Fields{
		"request_method"	: c.Request.Method,
		"request_uri"		: c.Request.RequestURI,
		"status_code"		: 0,
		"client_ip"			: c.ClientIP(),
		"latency_time"		: 0,
	})
}

