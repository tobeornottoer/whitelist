package model

type Admins struct {
	ID        	uint `gorm:"primarykey"`
	Account		string `gorm:"account"`
	Password	string `gorm:"password"`
}