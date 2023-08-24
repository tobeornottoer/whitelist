package model

import "gorm.io/gorm"

type TokenCost struct {
	gorm.Model
	EventID			uint64 		`gorm:"index"`
	AccountID		uint64 		`gorm:"index"`
	ActionType		int32		`gorm:"index"`
	Action 			string 		`gorm:"index"`
	ModelType		int32		`gorm:"index"`
	ModelName		string 		`gorm:"index"`
	Token			int32
}
