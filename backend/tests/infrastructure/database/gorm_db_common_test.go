package database_test

import (
	"os"
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/migration"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
)

var DbHandler *database.GormDbHandler
var UserDbHandler handler.UserDbHandler
var BlockDbHandler handler.BlockDbHandler
var MuteDbHandler handler.MuteDbHandler

func TestMain(m *testing.M) {
	// 前処理
	config := config.NewConfig(".env.test")
	DbHandler = database.NewGormDbHandler(config)
	UserDbHandler = database.NewUserDbHandler(DbHandler)
	BlockDbHandler = database.NewBlockDbHandler(DbHandler)
	MuteDbHandler = database.NewMuteHandler(DbHandler)

	// userテーブルから何も得られない場合はseederを実行
	user := &entity.User{}
	if err := UserDbHandler.Find(user, "id", "1"); err != nil || user.GetID() == 0 {
		migration.Seeder()
	}

	status := m.Run()

	os.Exit(status)
}
