package model

import (
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"gorm.io/gorm"
)

// for users table
type User struct {
	gorm.Model
	Name        string `gorm:"size:255;not null"`
	AccountName string `gorm:"size:255;not null"`
	TwitterID   string `gorm:"size:64;not null"`
	UserBlock   []UserBlock
	UserMute    []UserMute
}

func (u *User) FromUserDomain(userEntity *entity.User) {
	u.ID = userEntity.ID
	u.Name = userEntity.Name
	u.AccountName = userEntity.AccountName
	u.TwitterID = userEntity.TwitterID
}

func (u *User) ToUserDomain(i interface{}) {
	if s, ok := i.(*entity.User); ok {
		s.ID = u.ID
		s.Name = u.Name
		s.AccountName = u.AccountName
		s.TwitterID = u.TwitterID
	}
}
