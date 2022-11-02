package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

type GormDbUserHandler struct {
	//GormCommonDbHandler
	//GormDbEntityHandler[handler.DbHandler, entity.User, *model.UserModelForDomain[entity.User, model.User], model.User]
	GormDbEntityHandler[entity.User, model.User]
}

func NewUserDbHandler(gormDbHandler GormDbHandler) handler.UserDbHandler {
	userDbHandler := new(GormDbUserHandler)
	userDbHandler.db = gormDbHandler.Conn
	//userDbHandler.GormDbEntityHandler.db = userDbHandler.db
	userDbHandler.ModelForDomain = &model.UserModelForDomain{}

	return userDbHandler
}

/*
func (g *GormDbUserHandler) interfaceToEntity(userEntity interface{}) (*entity.User, error) {
	switch casted := userEntity.(type) {
	case *entity.User:
		return casted, nil
	default:
		return &entity.User{}, errors.New("interface is not entity.User")
	}

}
*/

// レコード1件をIDで取得
/*
func (g *GormDbUserHandler) First(userInterface interface{}, ID string) error {
	return g.GormDbEntityHandler.First(userInterface, ID)
}
*/

// レコードを作成
func (g *GormDbUserHandler) Create(userEntity interface{}) error {
	return g.GormDbEntityHandler.Create(userEntity)
	/*
	   userModel, err := g.entityToModel(userEntity)

	   	if err != nil {
	   		return err
	   	}

	   	if err = g.GormCommonDbHandler.Create(&userModel); err != nil {
	   		return err
	   	}

	   userModel.ToUserDomain(userEntity)

	   return err
	*/
}

// レコードをID基準で更新
func (g *GormDbUserHandler) Update(userEntity interface{}, ID string) error {
	return g.GormDbEntityHandler.Update(userEntity, ID)
	/*
	   userModel, err := g.entityToModel(userEntity)

	   	if err != nil {
	   		return err
	   	}

	   	if err = g.GormCommonDbHandler.Update(&userModel, ID); err != nil {
	   		return err
	   	}

	   userModel.ToUserDomain(&userEntity)

	   return err
	*/
}

// 条件に該当するレコードを新規作成または更新する
func (g *GormDbUserHandler) Upsert(userEntity interface{}, columnName string, searchValue string) error {
	return g.GormDbEntityHandler.Upsert(userEntity, columnName, searchValue)

	/*
		userModel, err := g.entityToModel(userEntity)
		if err != nil {
			return err
		}

		if err = g.GormCommonDbHandler.Upsert(&userModel, columnName, searchValue); err != nil {
			return err
		}

		userModel.ToUserDomain(userEntity)

		return err
	*/
}

// レコード1件を検索。※columnNameにユーザーからの入力値を絶対に使わないこと。
/*
func (g *GormDbUserHandler) Find(userEntity interface{}, columnName string, searchValue string) error {
	return g.GormDbEntityHandler.Find(userEntity, columnName, searchValue)
}
*/

/*
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
*/

func (g *GormDbUserHandler) FindAll(userEntity interface{}, columnName string, searchValue string) error {
	return g.GormDbEntityHandler.FindAll(userEntity, columnName, searchValue)
}

// ユーザーをTwitterIDで取得
func (g *GormDbUserHandler) FindByTwitterID(user interface{}, twitterID string) error {
	return g.Find(user, "twitter_id", twitterID)
}

// interfaceをentityにする。
/*
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
*/
