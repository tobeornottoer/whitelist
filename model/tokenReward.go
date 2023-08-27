package model

import (
	"gorm.io/gorm"
	"time"
)

type TokenReward struct {
	gorm.Model
	Token 		int
	Effect		time.Time
}