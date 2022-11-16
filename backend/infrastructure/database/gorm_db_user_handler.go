package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"gorm.io/gorm"
)

type GormDBUserHandler struct {
	GormDBEntityHandler[entity.User, model.User]
}

func NewUserDBHandler(conn handler.DBConnection) handler.UserDBHandler {
	userDBHandler := new(GormDBUserHandler)
	userDBHandler.db = conn.GetConnection().(*gorm.DB)
	userDBHandler.ModelForDomain = &model.UserModelForDomain{}

	return userDBHandler
}

// ユーザーをTwitterIDで取得
func (g *GormDBUserHandler) FindByTwitterID(user interface{}, twitterID string) error {
	return g.Find(user, "twitter_id", twitterID)
}
