package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"gorm.io/gorm"
)

type UserDbHandler struct {
	db *gorm.DB
}

func NewUserDbHandler(config *config.Config) gateway.DbHandler {
	gormHandler := NewGormHandler(config)

	userDbHandler := new(UserDbHandler)
	userDbHandler.db = gormHandler.Conn
	return userDbHandler
}

// ユーザー1件をIDで取得
func (u *UserDbHandler) First(user interface{}, userID string) error {
	if err := u.db.First(&user, userID).Error; err != nil {
		return err
    }
	return nil
}