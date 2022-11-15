package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormDbBlockHandler struct {
	GormDbEntityHandler[entity.Block, model.UserBlock]
}

func NewBlockDbHandler(conn handler.DbConnection) handler.BlockDbHandler {
	blockDbHandler := new(GormDbBlockHandler)
	blockDbHandler.db = conn.GetConnection().(*gorm.DB)
	blockDbHandler.ModelForDomain = &model.UserBlockModelForDomain{}

	return blockDbHandler
}

// ブロックをUserIDで全て取得
func (u *GormDbBlockHandler) FindAllByUserID(blockEntities interface{}, userID string) error {
	return u.FindAll(blockEntities, "user_id", userID)
}

// 新規ブロックを追加する（ただし追加済みのものも更新される）名称と実装が不一致
func (u *GormDbBlockHandler) CreateNewBlocks(recordSrc interface{}, columnName1 string, columnName2 string) error {
	blockEntities, err := u.GormDbEntityHandler.InterfaceToEntities(recordSrc)
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
