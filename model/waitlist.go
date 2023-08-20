package model

import "gorm.io/gorm"

type WaitList struct {
	gorm.Model
	Email 			string `gorm:"unique"`
	Code			string `gorm:"index"`
	Referral		uint32
	IP				string
	WhiteListFlag	bool 	`gorm:"index"`
	RegisterFlag	bool	`gorm:"index"`
	Unsubscribe		bool 	`gorm:"index"`
}
