package business

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"whitelist/model"
	"whitelist/utils"
)

func Dashboard(c *gin.Context){
	page,_		:= strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_	:= strconv.Atoi(c.DefaultQuery("pageSize","100"))
	searchUID, _ := strconv.Atoi(c.DefaultQuery("uid","0"))
	email		:= c.DefaultQuery("email","")
	sortField	:= c.DefaultQuery("sortField","")
	sortWay		:= c.DefaultQuery("sortWay","asc")

	var list []model.Account

	db	:= utils.GetDB()
	var total int64
	handle	:= db.Table("account")
	if searchUID > 0 {
		handle	= handle.Where("id = ?",searchUID)
	}
	if email != "" {
		handle 	= handle.Where("email = ?",email)
	}
	if sortField != "" {
		handle	= handle.Order(sortField + " " +sortWay)
	} else {
		handle	= handle.Order("last_token_cost_time desc,created_at desc,token_count desc")
	}
	result := handle.Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	if result.Error != nil {
		utils.CreateResponse(c).Json(http.StatusInternalServerError,result.Error.Error(),nil)
	} else {
		utils.CreateResponse(c).Success(gin.H{"list": list,"count":total})
	}
	return
}
