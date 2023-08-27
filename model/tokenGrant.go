package model

import "gorm.io/gorm"

type TokenGrant struct {
	gorm.Model
	AccountId   uint64 		`gorm:"index"`
	Date		string		`gorm:"index"`
	Token		uint64
	Way			int
}
