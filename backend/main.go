package main

import "github.com/faciam_dev/twitter_block2mute/backend/infrastructure"

func main() {
    db := infrastructure.NewDB()
    twitter := infrastructure.NewTwitter()
    r := infrastructure.NewRouting(db, twitter)
    r.Run()
}