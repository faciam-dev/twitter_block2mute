package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormDBBlockHandler struct {
	GormDBEntityHandler[entity.Block, model.UserBlock]
}

func NewBlockDBHandler(conn handler.DBConnection) handler.BlockDBHandler {
	blockDBHandler := new(GormDBBlockHandler)
	blockDBHandler.db = conn.GetConnection().(*gorm.DB)
	blockDBHandler.ModelForDomain = &model.UserBlockModelForDomain{}

	return blockDBHandler
}

// ブロックをUserIDで全て取得
func (u *GormDBBlockHandler) FindAllByUserID(blockEntities interface{}, userID string) error {
	return u.FindAll(blockEntities, "user_id", userID)
}

// 新規ブロックを追加する（ただし追加済みのものも更新される）名称と実装が不一致
func (u *GormDBBlockHandler) CreateNewBlocks(recordSrc interface{}, columnName1 string, columnName2 string) error {
	blockEntities, err := u.GormDBEntityHandler.InterfaceToEntities(recordSrc)
	if err != nil {
		return err
	}
	models := u.EntitiesToModels(*blockEntities)

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
