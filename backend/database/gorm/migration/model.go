package migration

import "gorm.io/gorm"

// for users
type User struct {
	gorm.Model
    Name string        `gorm:"size:255;not null"`
    AccountName string `gorm:"size:255;not null"`
    TwitterID string   `gorm:"size:64;not null"`
    UserBlock[] UserBlock
    UserMute[] UserMute
}

// for user_blocks
type UserBlock struct {
    gorm.Model
    UserID uint 
    TargetTwitterID string `gorm:"size:64;not null"`
    Flag int
    User User
}

// for user_mutes
type UserMute struct {
    gorm.Model
    UserID uint 
    TargetTwitterID string `gorm:"size:64;not null"`
    Flag int
    User User
}