package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
)

type GormDbUserHandler struct {
	GormCommonDbHandler
}

func NewUserDbHandler(gormDbHandler GormDbHandler) handler.UserDbHandler {
	userDbHandler := new(GormDbUserHandler)
	userDbHandler.db = gormDbHandler.Conn

	return userDbHandler
}

// ユーザーをTwitterIDで取得
func (u *GormDbUserHandler) FindByTwitterID(user interface{}, twitterID string) error {
	return u.Find(user, "twitter_id", twitterID)
}
