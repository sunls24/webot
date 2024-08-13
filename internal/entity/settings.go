package entity

import (
	"gorm.io/gorm"
)

type Settings struct {
	gorm.Model
	AvatarID string `gorm:"index"`
	Key      string
	Value    string
}
