package migration

import (
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
)

func Migration() {
    config := config.NewConfig()

    gormHandler := database.NewGormDbHandler(config)

	gormHandler.Connect().AutoMigrate(
		&User{},
		&UserBlock{},
		&UserMute{},
	)
}