package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"gorm.io/gorm/clause"
)

type GormDbBlockHandler struct {
	//GormCommonDbHandler
	GormDbEntityHandler[entity.Block, model.UserBlock]
}

func NewBlockDbHandler(gormDbHandler GormDbHandler) handler.BlockDbHandler {
	blockDbHandler := new(GormDbBlockHandler)
	blockDbHandler.db = gormDbHandler.Conn
	//blockDbHandler.GormDbEntityHandler.db = blockDbHandler.db
	blockDbHandler.ModelForDomain = &model.UserBlockModelForDomain{}

	return blockDbHandler
}

/*
// レコード1件をIDで取得
func (g *GormDbBlockHandler) First(blockEntity interface{}, ID string) error {
	blockModel, err := g.entityToModel(blockEntity)
	if err != nil {
		return err
	}

	if err = g.GormCommonDbHandler.First(&blockModel, ID); err != nil {
		return err
	}

	blockModel.ToDomain(blockEntity)
	return err
}

// レコードを作成
func (g *GormDbBlockHandler) Create(blockEntity interface{}) error {
	blockModel, err := g.entityToModel(blockEntity)
	if err != nil {
		return err
	}

	if err = g.GormCommonDbHandler.Create(&blockModel); err != nil {
		return err
	}

	blockModel.ToDomain(blockEntity)

	return err
}

// レコードをID基準で更新
func (g *GormDbBlockHandler) Update(blockEntity interface{}, ID string) error {
	blockModel, err := g.entityToModel(blockEntity)
	if err != nil {
		return err
	}

	if err = g.GormCommonDbHandler.Update(&blockModel, ID); err != nil {
		return err
	}

	blockModel.ToDomain(blockEntity)

	return err
}

// 条件に該当するレコードを新規作成または更新する
func (g *GormDbBlockHandler) Upsert(blockEntity interface{}, columnName string, searchValue string) error {
	blockModel, err := g.entityToModel(blockEntity)
	if err != nil {
		return err
	}

	if err = g.GormCommonDbHandler.Upsert(&blockModel, columnName, searchValue); err != nil {
		return err
	}

	blockModel.ToDomain(blockEntity)

	return err
}

// レコード1件を検索。※columnNameにユーザーからの入力値を絶対に使わないこと。
func (g *GormDbBlockHandler) Find(blockEntity interface{}, columnName string, searchValue string) error {
	blockModel, err := g.entityToModel(blockEntity)
	if err != nil {
		return err
	}

	if err = g.GormCommonDbHandler.Find(&blockModel, columnName, searchValue); err != nil {
		return err
	}

	blockModel.ToDomain(blockEntity)

	return err
}

// レコード全件を検索。※columnNameにユーザーからの入力値を絶対に使わないこと。
func (g *GormDbBlockHandler) FindAll(entities interface{}, columnName string, searchValue string) error {
	models, err := g.entitiesToModels(entities)
	if err != nil {
		return err
	}

	if err := g.GormCommonDbHandler.FindAll(models, columnName, searchValue); err != nil {
		return err
	}

	resultEntities := []entity.Block{}
	for _, model := range models {
		model.ToDomain(resultEntities[0])
	}

	return nil
}
*/

// ブロックをUserIDで全て取得
func (u *GormDbBlockHandler) FindAllByUserID(blockEntities interface{}, userID string) error {
	return u.FindAll(blockEntities, "user_id", userID)
}

// 記録済みブロックを変更せずに新規レコードを追加する
func (u *GormDbBlockHandler) CreateNewBlocks(recordSrc interface{}, columnName1 string, columnName2 string) error {
	blockEntities, err := u.GormDbEntityHandler.InterfaceToEntities(recordSrc)
	if err != nil {
		return err
	}
	models := u.EntitiesToModels(*blockEntities)

	err = u.db.Debug().Clauses(clause.OnConflict{
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

/*
// interface(entity)をmodelにする。
func (u *GormDbBlockHandler) entityToModel(blockEntity interface{}) (model.UserBlock, error) {
	blockModel := model.UserBlock{}

	switch casted := blockEntity.(type) {
	case *entity.Block:
		blockModel.FromDomain(casted)
	default:
		return blockModel, errors.New("argument type is not entity.Block")
	}
	return blockModel, nil
}

// interface([]]entity)を[]modelにする。
func (u *GormDbBlockHandler) entitiesToModels(blockEntity interface{}) ([]model.UserBlock, error) {
	blockModels := []model.UserBlock{}

	switch casted := blockEntity.(type) {
	case *[]entity.Block:
		for _, entity := range *casted {
			model := model.UserBlock{}
			model.FromDomain(&entity)
			blockModels = append(blockModels, model)
		}
	default:
		return blockModels, errors.New("argument type is not []entity.Block")
	}
	return blockModels, nil

}
*/
