package business

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"time"
	"whitelist/utils"
)

type AccountEvent struct{
	Event 	string 		`json:"event"`
	Time	time.Time	`json:"time"`
	Status	string 		`json:"status"`
	Latency decimal.Decimal	`json:"latency"`
	Tokens	uint		`json:"tokens"`
}

func GetEventHistory(uid uint64, page, pageSize int, dateStart, dateEnd time.Time) (gin.H,error) {
	var events []AccountEvent
	db	:= utils.GetDB()
	var total	int64
	result := db.Table("token_costs").Select("events.event,events.created_at AS time,events.status,events.latency,token_costs.token AS tokens").Joins("inner join events on events.id = token_costs.event_id").Where("token_costs.account_id = ?",uid).Where("events.created_at between ? and ?",dateStart,dateEnd).Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Order("events.id desc").Scan(&events)

	if result.Error != nil {
		return gin.H{},result.Error
	}

	return gin.H{
		"list": events,"count":total,
	},nil

}

type EventSummary struct{
	Event 		string 		`json:"event"`
	ModelType	string		`json:"model_type"`
	Numbers		string 		`json:"numbers"`
	Tokens		uint		`json:"tokens"`
	Percentage 	string		`json:"percentage"`
}

type EventSummarySum struct {
	Total		uint
}

func DurationSummary(uid uint64, dateStart, dateEnd time.Time) (gin.H,error) {
	var eventSummary []EventSummary
	var total EventSummarySum
	db	:= utils.GetDB()
	tr 	:= db.Table("token_costs").Select("sum(token_costs.token) as total").Joins("inner join events on events.id = token_costs.event_id").Where("token_costs.account_id = ?",uid).Where("events.created_at between ? and ?",dateStart,dateEnd).Scan(&total)

	if tr.Error != nil {
		return gin.H{},tr.Error
	}

	esr	:= db.Table("token_costs").Select("events.event,token_costs.model_type,count(*) AS numbers,token_costs.token AS tokens").Joins("inner join events on events.id = token_costs.event_id").Where("token_costs.account_id = ?",uid).Where("events.created_at between ? and ?",dateStart,dateEnd).Group("events.event").Scan(&eventSummary)

	if esr.Error != nil {
		return gin.H{},esr.Error
	}

	for _,es	:= range eventSummary {
		es.Percentage	= fmt.Sprintf("%.2f%%",float64(es.Tokens / total.Total) * 100)
	}

	return gin.H{
		"list":eventSummary,"total":total.Total,
	},nil
}
