package database

import (
	"errors"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

type GormDbUserHandler struct {
	GormCommonDbHandler
}

func NewUserDbHandler(gormDbHandler GormDbHandler) handler.UserDbHandler {
	userDbHandler := new(GormDbUserHandler)
	userDbHandler.db = gormDbHandler.Conn

	return userDbHandler
}

// レコード1件をIDで取得
func (g *GormDbUserHandler) First(userEntity interface{}, ID string) error {
	userModel, err := g.entityToModel(userEntity)
	if err != nil {
		return err
	}

	if err = g.GormCommonDbHandler.First(&userModel, ID); err != nil {
		return err
	}

	userModel.ToUserDomain(userEntity)

	return err
}

// レコードを作成
func (g *GormDbUserHandler) Create(userEntity interface{}) error {
	userModel, err := g.entityToModel(userEntity)
	if err != nil {
		return err
	}

	if err = g.GormCommonDbHandler.Create(&userModel); err != nil {
		return err
	}

	userModel.ToUserDomain(userEntity)

	return err
}

// レコードをID基準で更新
func (g *GormDbUserHandler) Update(userEntity interface{}, ID string) error {
	userModel, err := g.entityToModel(userEntity)
	if err != nil {
		return err
	}

	if err = g.GormCommonDbHandler.Update(&userModel, ID); err != nil {
		return err
	}

	userModel.ToUserDomain(userEntity)

	return err
}

// 条件に該当するレコードを新規作成または更新する
func (g *GormDbUserHandler) Upsert(userEntity interface{}, columnName string, searchValue string) error {
	userModel, err := g.entityToModel(userEntity)
	if err != nil {
		return err
	}

	if err = g.GormCommonDbHandler.Upsert(&userModel, columnName, searchValue); err != nil {
		return err
	}

	userModel.ToUserDomain(userEntity)

	return err
}

// レコード1件を検索。※columnNameにユーザーからの入力値を絶対に使わないこと。
func (g *GormDbUserHandler) Find(userEntity interface{}, columnName string, searchValue string) error {
	userModel, err := g.entityToModel(userEntity)
	if err != nil {
		return err
	}

	if err = g.GormCommonDbHandler.Find(&userModel, columnName, searchValue); err != nil {
		return err
	}

	userModel.ToUserDomain(userEntity)

	return err
}

// ユーザーをTwitterIDで取得
func (u *GormDbUserHandler) FindByTwitterID(user interface{}, twitterID string) error {
	return u.Find(user, "twitter_id", twitterID)
}

// interfaceをentityにする。
func (u *GormDbUserHandler) entityToModel(userEntity interface{}) (model.User, error) {
	userModel := model.User{}

	switch casted := userEntity.(type) {
	case *entity.User:
		userModel.FromUserDomain(casted)
	default:
		return userModel, errors.New("argment type is not entity.User")
	}
	return userModel, nil
}
