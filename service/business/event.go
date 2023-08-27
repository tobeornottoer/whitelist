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
	tx		:= db.Table("token_costs")
	tx		=  tx.Select("events.event,events.created_at AS time,events.status,events.latency,token_costs.token AS tokens")
	tx		=  tx.Joins("inner join events on events.id = token_costs.event_id").Where("token_costs.account_id = ?",uid)
	tx		=  tx.Where("events.created_at between ? and ?",dateStart,dateEnd)
	tx		=  tx.Count(&total)
	result :=  tx.Offset((page - 1) * pageSize).Limit(pageSize).Order("events.id desc").Scan(&events)

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
	tx	:= db.Table("token_costs")
	tx	=  tx.Select("sum(token_costs.token) as total")
	tx	=  tx.Joins("inner join events on events.id = token_costs.event_id")
	tx	=  tx.Where("token_costs.account_id = ?",uid)
	tx	=  tx.Where("events.created_at between ? and ?",dateStart,dateEnd)
	tr  := tx.Scan(&total)

	if tr.Error != nil {
		return gin.H{},tr.Error
	}

	ty	:= db.Table("token_costs")
	ty	=  ty.Select("any_value(events.event) as event,any_value(token_costs.model_type) as model_type,count(*) AS numbers,any_value(token_costs.token) AS tokens")
	ty	=  ty.Joins("inner join events on events.id = token_costs.event_id")
	ty	=  ty.Where("token_costs.account_id = ?",uid)
	ty	=  ty.Where("events.created_at between ? and ?",dateStart,dateEnd)
	esr	:= ty.Group("events.event").Scan(&eventSummary)

	if esr.Error != nil {
		return gin.H{},esr.Error
	}

	for i,es	:= range eventSummary {
		es.Percentage	= fmt.Sprintf("%.2f%%",float64(es.Tokens / total.Total) * 100)
		eventSummary[i] = es
	}

	return gin.H{
		"list":eventSummary,"total":total.Total,
	},nil
}
