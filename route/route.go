package route

import (
	"github.com/gin-gonic/gin"
	"whitelist/middleware"
	"whitelist/service/admins"
	whitelistService "whitelist/service/whitelist"
)


func Set(r *gin.Engine) {
	whitelist(r)
	login(r)
	initAdmin(r)
}

func login(r *gin.Engine){
	r.POST("/login",func(c *gin.Context){
		admins.Login(c)
	})
}

func initAdmin(r *gin.Engine){
	r.GET("/init/admin",func(c *gin.Context){
		admins.Register(c)
	})
}

// 白名单路由组
func whitelist(r *gin.Engine){
	authorized := r.Group("/whitelist")
	authorized.Use(middleware.UserAuthorized())
	{
		authorized.GET("",func(c *gin.Context){
			whitelistService.Get(c)
		})
		authorized.POST("/approve",func(c *gin.Context){
			whitelistService.Approve(c)
		})
		authorized.POST("/add/email",func(c *gin.Context){
			whitelistService.AddEmail(c)
		})
		authorized.POST("/upload/excel",func(c *gin.Context){
			whitelistService.UploadEmailExcel(c)
		})
		authorized.POST("/import/excel",func(c *gin.Context){
			whitelistService.ImportEmails(c)
		})
	}
}