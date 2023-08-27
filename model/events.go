package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Events struct{
	gorm.Model
	Event 				string 			`json:"event"`
	Status				string 			`json:"status"`
	Latency				decimal.Decimal	`json:"latency"`
	PromptTokens		uint			`json:"prompt_tokens"`
	CompletionTokens	uint			`json:"completion_tokens"`
	Cost				decimal.Decimal	`json:"cost"`
}
