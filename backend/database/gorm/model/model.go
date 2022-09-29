package model

import "gorm.io/gorm"

// for user_blocks
type UserBlock struct {
	gorm.Model
	UserID          uint
	TargetTwitterID string `gorm:"size:64;not null"`
	Flag            int
	User            User
}

// for user_mutes
type UserMute struct {
	gorm.Model
	UserID          uint
	TargetTwitterID string `gorm:"size:64;not null"`
	Flag            int
	User            User
}
