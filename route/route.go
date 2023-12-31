package route

import (
	"github.com/gin-gonic/gin"
	"whitelist/middleware"
	"whitelist/service/admins"
	business2 "whitelist/service/business"
	whitelistService "whitelist/service/whitelist"
)


func Set(r *gin.Engine) {
	whitelist(r)
	login(r)
	initAdmin(r)
	business(r)
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

func business(r *gin.Engine) {
	authorized := r.Group("")
	authorized.Use(middleware.UserAuthorized())
	{
		authorized.GET("/token/dashboard",func(c *gin.Context){
			business2.Dashboard(c)
		})
		authorized.GET("/token/info",func(c *gin.Context){
			business2.Info(c)
		})
		authorized.POST("/update/token",func(c *gin.Context){
			business2.UpdateToken(c)
		})
		authorized.GET("/token/grant",func(c *gin.Context){
			business2.GetTokenGrantLogs(c)
		})
		authorized.POST("/update/token/reward",func(c *gin.Context){
			business2.UpdateDefaultTokenReward(c)
		})
		authorized.POST("/import/token/grant",func(c *gin.Context){
			business2.BatchGrant(c)
		})
		authorized.GET("/module/usage",func(c *gin.Context){
			business2.ModelUsage(c)
		})
		authorized.GET("/module/usage/:model_type",func(c *gin.Context){
			business2.ModelUsageByType(c)
		})
	}
}