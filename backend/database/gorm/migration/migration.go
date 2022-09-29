package migration

import (
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
)

func Migration() {
	config := config.NewConfig(".env")

	gormHandler := database.NewGormDbHandler(config)

	gormHandler.Connect().AutoMigrate(
		&model.User{},
		&model.UserBlock{},
		&model.UserMute{},
	)
}
