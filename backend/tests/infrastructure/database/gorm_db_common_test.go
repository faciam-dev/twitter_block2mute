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

var DBHandler handler.DBHandler
var UserDBHandler handler.UserDBHandler
var BlockDBHandler handler.BlockDBHandler
var MuteDBHandler handler.MuteDBHandler

func TestMain(m *testing.M) {
	// 前処理
	config := config.NewConfig(".env.test")

	dbConnection := database.NewGormDBConnectionByConfig(config)
	DBHandler = database.NewGormDBHandler(dbConnection)

	UserDBHandler = database.NewUserDBHandler(DBHandler.Connect())
	BlockDBHandler = database.NewBlockDBHandler(DBHandler.Connect())
	MuteDBHandler = database.NewMuteHandler(DBHandler.Connect())

	// userテーブルから何も得られない場合はseederを実行
	user := &entity.User{}
	if err := UserDBHandler.Find(user, "id", "1"); err != nil || user.GetID() == 0 {
		migration.Seeder()
	}

	status := m.Run()

	os.Exit(status)
}
