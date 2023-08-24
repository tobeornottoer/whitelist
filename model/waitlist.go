package model

import "gorm.io/gorm"

type WaitList struct {
	gorm.Model
	Email 			string 	`gorm:"unique" json:"email"`
	Code			string 	`gorm:"index" json:"code"`
	Referral		uint32 	`json:"referral"`
	IP				string	`json:"ip"`
	WhiteListFlag	bool 	`gorm:"index" json:"white_list_flag"`
	RegisterFlag	bool	`gorm:"index" json:"register_flag"`
	Unsubscribe		bool 	`gorm:"index" json:"unsubscribe"`
}
