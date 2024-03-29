package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormDBMuteHandler struct {
	GormDBEntityHandler[entity.Mute, model.UserMute]
}

func NewMuteHandler(conn handler.DBConnection) handler.MuteDBHandler {
	dbHandler := new(GormDBMuteHandler)
	dbHandler.db = conn.GetConnection().(*gorm.DB)
	dbHandler.ModelForDomain = &model.UserMuteModelForDomain{}

	return dbHandler
}

// ミュートをUserIDで全て取得
func (u *GormDBMuteHandler) FindAllByUserID(entities interface{}, userID string) error {
	return u.FindAll(entities, "user_id", userID)
}

// 新規ミュートを追加する（ただし追加済みのものも更新される）名称と実装が不一致
func (u *GormDBMuteHandler) CreateNew(recordSrc interface{}, columnName1 string, columnName2 string) error {
	entities, err := u.GormDBEntityHandler.InterfaceToEntities(recordSrc)
	if err != nil {
		return err
	}
	models := u.EntitiesToModels(*entities)

	err = u.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "target_twitter_id"}},
		//Where: clause.Where{Exprs: []clause.Expression{clause.Eq{Colunm: columnName1,  columnName2}}},
		//Where:     clause.Where{Exprs: []clause.Expression{clause.Eq{Column: columnName, Value: searchValue}}},
		UpdateAll: true,
	}).Create(&models).Error

	if err != nil {
		return err
	}
	return nil
}
