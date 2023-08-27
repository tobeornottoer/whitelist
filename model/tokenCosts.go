package model

import "gorm.io/gorm"

type TokenCosts struct {
	gorm.Model
	EventID			uint 		`gorm:"index"`
	AccountID		uint 		`gorm:"index"`
	ActionType		int32		`gorm:"index"`
	Action 			string 		`gorm:"index"`
	ModelType		int32		`gorm:"index"`
	ModelName		string 		`gorm:"index"`
	Token			int32
}
