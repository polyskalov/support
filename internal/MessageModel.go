package internal

import (
	"gorm.io/gorm"
	"time"
)

type MessageModel struct {
	gorm.Model
	ID         string `gorm:"type:uuid;default:gen_random_uuid()"`
	ExternalId int
	Time       time.Time `gorm:"index"`
	User       User
	Text       string
}
