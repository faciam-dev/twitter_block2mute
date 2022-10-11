package model

import "gorm.io/gorm"

// for user_mutes
type UserMute struct {
	gorm.Model
	UserID          uint   `gorm:"uniqueIndex:idx_user_id_target_twitter_id"`
	TargetTwitterID string `gorm:"uniqueIndex:idx_user_id_target_twitter_id;size:64;not null"`
	Flag            int
	User            User
}
