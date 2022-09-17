package main

import (
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure"
)

func main() {
    twitter := infrastructure.NewTwitter()
    r := infrastructure.NewRouting(twitter)
    r.Run()
}