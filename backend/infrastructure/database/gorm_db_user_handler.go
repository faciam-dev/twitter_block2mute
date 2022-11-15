package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"gorm.io/gorm"
)

type GormDbUserHandler struct {
	GormDbEntityHandler[entity.User, model.User]
}

func NewUserDbHandler(conn handler.DbConnection) handler.UserDbHandler {
	userDbHandler := new(GormDbUserHandler)
	userDbHandler.db = conn.GetConnection().(*gorm.DB)
	userDbHandler.ModelForDomain = &model.UserModelForDomain{}

	return userDbHandler
}

// ユーザーをTwitterIDで取得
func (g *GormDbUserHandler) FindByTwitterID(user interface{}, twitterID string) error {
	return g.Find(user, "twitter_id", twitterID)
}
