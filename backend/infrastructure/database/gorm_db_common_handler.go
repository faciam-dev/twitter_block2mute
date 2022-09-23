package database

import (
	"errors"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormCommonDbHandler struct {
	db *gorm.DB
}

func NewGormCommonDbHandler(databaseDbHandler GormDbHandler) handler.DbHandler {
	gormCommonDbHandler := new(GormCommonDbHandler)
	gormCommonDbHandler.db = databaseDbHandler.Conn
	return gormCommonDbHandler
}

// レコード1件をIDで取得
func (g *GormCommonDbHandler) First(model interface{}, ID string) error {
	if err := g.db.First(&model, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

// レコードを作成
func (g *GormCommonDbHandler) Create(recordSrc interface{}) error {
	if err := g.db.Create(recordSrc).Error; err != nil {
		return err
	}
	return nil
}

// レコードをID基準で更新
func (g *GormCommonDbHandler) Update(recordSrc interface{}, ID string) error {
	if err := g.db.Model(&recordSrc).Where("id = ?", ID).Updates(recordSrc).Error; err != nil {
		return err
	}
	return nil
}

// レコード1件を検索。※columnNameにユーザーからの入力値を絶対に使わないこと。
func (g *GormCommonDbHandler) Find(model interface{}, columnName string, searchValue string) error {
	/*
		if err := g.db.Raw("SELECT * FROM `users` WHERE twitter_id = ?", searchValue).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
	*/
	//SELECT * FROM `users` WHERE twitter_id = '1234567890'
	if err := g.db.Where(columnName+" = ?", searchValue).Find(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

// 条件に該当するレコードを新規作成または更新する
func (g *GormCommonDbHandler) Upsert(recordSrc interface{}, columnName string, searchValue string) error {
	err := g.db.Clauses(clause.OnConflict{
		//Columns: []clause.Column{{Name: columnName}},
		Where:     clause.Where{Exprs: []clause.Expression{clause.Eq{Column: columnName, Value: searchValue}}},
		UpdateAll: true,
	}).Create(recordSrc).Error

	if err != nil {
		return err
	}
	return nil
}
