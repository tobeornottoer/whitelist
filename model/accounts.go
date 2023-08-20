package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Username	string		`gorm:"index;null;unique"`
	Password	[]byte
	Nickname	string 		`gorm:"null"`
	AuthingID	string 		`gorm:"null;unique;index"`
	TokenCount	uint64
}
