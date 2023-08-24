package model

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	gorm.Model
	Username	string		`gorm:"index;null;unique" json:"username"`
	Password	[]byte		`json:"-"`
	Nickname	string 		`gorm:"null" json:"nickname"`
	Email		string 		`gorm:"unique" json:"email"`
	AuthingID	string 		`gorm:"null;unique;index" json:"authing_id"`
	TokenCount	uint64		`json:"token_count"`
	LastTokenCostTime	time.Time	`json:"last_token_cost_time"`
}
