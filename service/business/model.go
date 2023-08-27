package business

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
	"whitelist/utils"
)

type Summary struct{
	ModelType 		[]string 	`json:"model_type"`
	TotalTokens		uint		`json:"total_tokens"`
	TotalEvent		int			`json:"total_event"`
	TotalCosts		decimal.Decimal	`json:"total_costs"`
}

func ModelUsage(c *gin.Context) {
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
	var types []string
	db		:= utils.GetDB()
	result 	:= db.Table("token_costs").Where("created_at between ? and ?",dateStart,dateEnd).Group("model_name").Pluck("model_type",&types)
	if result.Error != nil {
		utils.CreateResponse(c).ServerError("Error querying model type")
		return
	}
	summary,err 	:= GetSummaryByModelType(types[0],dateStart,dateEnd)
	if err != nil {
		utils.CreateResponse(c).ServerError(err.Error())
		return
	}
	logs,err 		:= GetModelEventLogs(types[0],dateStart,dateEnd,page,pageSize)
	if err != nil {
		utils.CreateResponse(c).ServerError(err.Error())
		return
	}
	utils.CreateResponse(c).Success(gin.H{
		"model_type":types,
		"summary": summary,
		"logs": logs,
	})
	return
}

func ModelUsageByType(c *gin.Context){
	start		:= c.DefaultQuery("start","")
	end			:= c.DefaultQuery("end","")
	page,_		:= strconv.Atoi(c.DefaultQuery("page","0"))
	pageSize,_	:= strconv.Atoi(c.DefaultQuery("pageSize","100"))
	if start == "" || end == "" {
		utils.CreateResponse(c).BadRequest("Please select a date range")
		return
	}
	dateStart,err	:= time.Parse(start,"2006-01-02 15:04:05")
	if err != nil {
		utils.CreateResponse(c).BadRequest("Incorrect time range")
		return
	}
	dateEnd,err		:= time.Parse(end,"2006-01-02 15:04:05")
	if err != nil {
		utils.CreateResponse(c).BadRequest("Incorrect time range")
		return
	}
	modelType 	:= c.Param("model_type")
	summary,err 	:= GetSummaryByModelType(modelType,dateStart,dateEnd)
	if err != nil {
		utils.CreateResponse(c).ServerError(err.Error())
		return
	}
	logs,err 		:= GetModelEventLogs(modelType,dateStart,dateEnd,page,pageSize)
	if err != nil {
		utils.CreateResponse(c).ServerError(err.Error())
		return
	}
	utils.CreateResponse(c).Success(gin.H{
		"summary": summary,
		"logs": logs,
	})
	return
}

type ModelEventLogs struct {
	Event 		string 			`json:"event"`
	Time 		time.Time		`json:"time"`
	Status		string			`json:"status"`
	Latency 	decimal.Decimal	`json:"latency"`
	Prompt		uint			`json:"prompt"`
	Completion 	uint			`json:"completion"`
	Cost		decimal.Decimal	`json:"cost"`
}

type ModelEventSummary struct {
	TotalTokens		uint		`json:"total_tokens"`
	TotalEvent		int			`json:"total_event"`
	TotalCosts		decimal.Decimal	`json:"total_costs"`
}

func GetSummaryByModelType(modelType string, start, end time.Time) (ModelEventSummary,error){
	var	summary ModelEventSummary
	db		:= utils.GetDB()
	tx		:= db.Table("token_costs").Joins("events on events.id = token_costs.event_id")
	tx		=  tx.Where("token_costs.model_name = ?",modelType)
	tx		=  tx.Where("token_costs.created_at between ? and ?",start,end)
	r		:=  tx.Select("SUM(events.prompt_tokens + events.completion_tokens) AS total_tokens,count(*) AS total_event,SUM(events.cost) AS total_costs").First(&summary)
	return summary,r.Error
}

func GetModelEventLogs(modelType string, start, end time.Time,page,pageSize int) ([]ModelEventLogs,error) {
	var list []ModelEventLogs
	db		:= utils.GetDB()
	tx		:= db.Table("token_costs").Joins("events on events.id = token_costs.event_id")
	tx		=  tx.Where("token_costs.model_name = ?",modelType)
	tx		=  tx.Where("token_costs.created_at between ? and ?",start,end)
	tx		=  tx.Select("events.event,events.created_at as `time`,events.status,events.latency,events.prompt_tokens as prompt,events.completion_tokens,events.cost")
	r		:= tx.Offset((page - 1) * pageSize).Limit(pageSize).Order("events.id desc").Scan(&list)
	return list,r.Error
}