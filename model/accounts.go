package model

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	gorm.Model
	Username	string		`gorm:"index;null;unique"`
	Password	[]byte
	Nickname	string 		`gorm:"null"`
	AuthingID	string 		`gorm:"null;unique;index"`
	TokenCount	uint64
	LastTokenCostTime	time.Time
}
