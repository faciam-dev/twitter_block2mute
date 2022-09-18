package main

import (
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/twitterapi"
)

func main() {
    config := config.NewConfig()

    dbHandler := database.NewUserDbHandler(config)
    anaconderHandler := twitterapi.NewAnaconderHandler(config)
    //sessionHandler := framework.NewGinSessionHandler(NewGinSessionHandler(config)

    r := infrastructure.NewRouting(config, dbHandler, anaconderHandler)
    r.Run()
}