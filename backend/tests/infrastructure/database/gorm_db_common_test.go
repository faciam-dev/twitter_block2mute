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

var UserDbHandler handler.UserDbHandler
var BlockDbHandler handler.BlockDbHandler
var MuteDbHandler handler.MuteDbHandler

func TestMain(m *testing.M) {
	// 前処理
	config := config.NewConfig(".env.test")

	dbHandler := database.NewGormDbHandler(config)

	UserDbHandler = database.NewUserDbHandler(dbHandler)
	BlockDbHandler = database.NewBlockDbHandler(dbHandler)
	MuteDbHandler = database.NewMuteHandler(dbHandler)

	// userテーブルから何も得られない場合はseederを実行
	user := &entity.User{}
	if err := UserDbHandler.Find(user, "id", "1"); err != nil || user.ID == 0 {
		migration.Seeder()
	}

	// トランザクション
	UserDbHandler.Begin()
	BlockDbHandler.Begin()
	MuteDbHandler.Begin()

	status := m.Run()

	// 後処理
	// トランザクションを戻す
	// MySQLではauto_incrementは戻らない
	UserDbHandler.Rollback()
	BlockDbHandler.Rollback()
	MuteDbHandler.Rollback()

	os.Exit(status)
}
