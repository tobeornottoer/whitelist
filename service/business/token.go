package business

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
	"whitelist/logger"
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
		utils.CreateResponse(c).ServerError(result.Error.Error())
	} else {
		utils.CreateResponse(c).Success(gin.H{"list": list,"count":total})
	}
	return
}

func Info(c *gin.Context) {
	uid,_		:= strconv.ParseUint(c.DefaultQuery("uid","0"),10,64)
	start		:= c.DefaultQuery("start","")
	end			:= c.DefaultQuery("end","")
	page,_		:= strconv.Atoi(c.DefaultQuery("page","0"))
	pageSize,_	:= strconv.Atoi(c.DefaultQuery("pageSize","100"))
	if start == "" || end == "" {
		utils.CreateResponse(c).BadRequest("Please select a date range")
		return
	}
	dateStart,err	:= time.Parse("2006-01-02 15:04:05",start)
	if err != nil {
		utils.CreateResponse(c).BadRequest("Incorrect time range")
		return
	}
	dateEnd,err		:= time.Parse("2006-01-02 15:04:05",end)
	if err != nil {
		utils.CreateResponse(c).BadRequest("Incorrect time range")
		return
	}
	if uid <= 0 {
		utils.CreateResponse(c).BadRequest("parameter is incorrect")
		return
	}
	db 		:= utils.GetDB()
	var account model.Account
	fc		:= db.Model(&model.Account{}).Where("id = ?" ,uid).First(&account)
	if errors.Is(fc.Error,gorm.ErrRecordNotFound) {
		utils.CreateResponse(c).BadRequest("Unable to find the account")
		return
	}

	summary,sErr	:= DurationSummary(uid,dateStart,dateEnd)
	if sErr != nil {
		logger.RuntimeLog.Error(sErr.Error())
	}
	history,hErr	:= GetEventHistory(uid,page,pageSize,dateStart,dateEnd)
	if hErr != nil {
		logger.RuntimeLog.Error(hErr.Error())
	}

	utils.CreateResponse(c).Success(gin.H{
		"info":account,
		"summary":summary,
		"history":history,
	})
}


func UpdateToken(c *gin.Context) {
	uid,uidErr		:= strconv.ParseUint(c.PostForm("uid"),10,64)
	tokens,tokenErr	:= strconv.ParseUint(c.PostForm("tokens"),10,64)
	if uid <= 0 || uidErr != nil || tokenErr != nil {
		utils.CreateResponse(c).BadRequest("parameter is incorrect")
		return
	}
	db	:= utils.GetDB()
	result := db.Model(&model.Account{}).Where("id = ?", uid).Update("token_count", tokens)
	if result.Error != nil {
		utils.CreateResponse(c).ServerError("Modification failed")
		return
	}
	utils.CreateResponse(c).Success(gin.H{})
	return
}

func GetTokenGrantLogs(c *gin.Context){
	type grantList struct {
		Date		string 	`json:"date"`
		Number		int		`json:"number"`
		TotalTokens	uint	`json:"total_tokens"`
	}
	page,_		:= strconv.Atoi(c.DefaultQuery("page","0"))
	pageSize,_	:= strconv.Atoi(c.DefaultQuery("pageSize","100"))
	var list []grantList
	var total int64
	db		:= utils.GetDB()
	result	:= db.Model(&model.TokenGrant{}).Select("date,count(*) as number,sum(token) as total_tokens").Group("date").Count(&total).Order("date desc").Offset((page - 1) * pageSize).Limit(pageSize).Scan(&list)
	if result.Error != nil {
		utils.CreateResponse(c).ServerError(result.Error.Error())
		return
	}
	utils.CreateResponse(c).Success(gin.H{
		"tokens":GetDefaultTokenReward(),
		"list":list,
		"count":total,
	})
	return
}

func GetDefaultTokenReward() int {
	db	:= utils.GetDB()
	var reward model.TokenReward
	result	:= db.Model(&model.TokenReward{}).Order("id desc").First(&reward)
	if errors.Is(result.Error,gorm.ErrRecordNotFound) {
		rt,err 	:= strconv.Atoi(os.Getenv("SYSTEM_DEFAULT_TOKEN_REWARD"))
		if err != nil {
			return 0
		}
		return rt
	}
	return reward.Token
}

func UpdateDefaultTokenReward(c *gin.Context){
	token,err 	:= strconv.Atoi(c.PostForm("tokens"))
	if err != nil {
		utils.CreateResponse(c).BadRequest("parameter is incorrect")
		return
	}
	tomorrow,_ 	:= time.Parse("2006-01-02",time.Now().Format("2006-01-02"))
	effect		:= tomorrow.AddDate(0,0,1)
	reward		:= model.TokenReward{
		Token: token,Effect: effect,
	}
	db			:= utils.GetDB()
	result		:= db.Create(&reward)
	if result.Error != nil {
		utils.CreateResponse(c).ServerError("Update failed")
		return
	}
	utils.CreateResponse(c).Success(gin.H{})
	return
}

func BatchGrant(c *gin.Context){
	code 	:= c.DefaultPostForm("code","")
	if code == "" {
		utils.CreateResponse(c).BadRequest("file not found")
		return
	}
	fileByte,err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		utils.CreateResponse(c).BadRequest("file not found")
		return
	}
	file		:= string(fileByte)
	rootPath,_	:= os.Getwd()
	filePath	:= rootPath + file
	importErr 	:= utils.BatchGrant(filePath)
	_			= os.Remove(filePath)
	if importErr != nil {
		utils.CreateResponse(c).ServerError(importErr.Error())
		return
	}
	utils.CreateResponse(c).Success(gin.H{})
	return
}