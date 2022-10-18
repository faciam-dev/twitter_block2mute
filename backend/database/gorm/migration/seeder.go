package migration

import (
	"log"

	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
)

var users = []model.User{
	{
		Name:        "name1",
		AccountName: "Name1",
		TwitterID:   "1234567890",
	},
}

var userBlocks = []model.UserBlock{
	{
		UserID:          1,
		TargetTwitterID: "12345678901",
		Flag:            1,
	},
	{
		UserID:          1,
		TargetTwitterID: "12345678902",
		Flag:            1,
	},
	{
		UserID:          1,
		TargetTwitterID: "12345678903",
		Flag:            0,
	},
}

var userMutes = []model.UserMute{
	{
		UserID:          1,
		TargetTwitterID: "12345678904",
		Flag:            1,
	},
	{
		UserID:          1,
		TargetTwitterID: "12345678905",
		Flag:            0,
	},
}

func Seeder() {
	config := config.NewConfig(".env.test")

	gormHandler := database.NewGormDbHandler(config)

	err := gormHandler.Conn.Debug().Migrator().DropTable(
		&model.User{},
		&model.UserBlock{},
		&model.UserMute{},
	)
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = gormHandler.Conn.Debug().AutoMigrate(
		&model.User{},
		&model.UserBlock{},
		&model.UserMute{},
	)
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i := range users {
		err = gormHandler.Conn.Debug().Model(&model.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

	err = gormHandler.Conn.Model(&model.UserBlock{}).Create(&userBlocks).Error
	if err != nil {
		log.Fatalf("cannot seed user_blocks table: %v", err)
	}

	err = gormHandler.Conn.Model(&model.UserMute{}).Create(&userMutes).Error
	if err != nil {
		log.Fatalf("cannot seed user_mutes table: %v", err)
	}
}
