package main

import (
	"flag"

	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/migration"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/framework"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/logger"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/twitterapi"
)

func main() {
	mode := flag.String("mode", "", "mode")
	flag.Parse()

	if *mode == "migration" {
		migration.Migration()
	} else {
		server()
	}
}

func server() {
	config := config.NewConfig(".env")

	loggerHandler := logger.NewZapHandler(config)
	dbConnection := database.NewGormDBConnectionByConfig(config)
	twitterHandler := twitterapi.NewGotwiHandler(config)

	r := framework.NewRouting(config, loggerHandler, dbConnection, twitterHandler)
	r.Run()
}
