package utils

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strconv"
	"whitelist/logger"
)

var _db *gorm.DB

func init() {
	host 	:= os.Getenv("MYSQL_HOST")
	port,_	:= strconv.Atoi(os.Getenv("MYSQL_PORT"))
	user	:= os.Getenv("MYSQL_USERNAME")
	password:= os.Getenv("MYSQL_PASSWORD")
	database:= os.Getenv("MYSQL_DATABASE")
	timeout := os.Getenv("MYSQL_TIMEOUT")
	maxPool,_ := strconv.Atoi(os.Getenv("MYSQL_MAX_POOL"))
	maxIdle,_ := strconv.Atoi(os.Getenv("MYSQL_MAX_IDLE"))

	if maxPool == 0 {
		maxPool = 50
	}
	if maxIdle == 0 {
		maxIdle = 20
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%ss",
		user,
		password, 
		host, 
		port,
		database,
		timeout)
	fmt.Println(dsn)
	var err error

	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.RuntimeLog.Error("连接数据库失败, error=" + err.Error())
	}

	sqlDB, _ := _db.DB()

	sqlDB.SetMaxOpenConns(maxPool)
	sqlDB.SetMaxIdleConns(maxIdle)
}

func GetDB() *gorm.DB {
	return _db
}
