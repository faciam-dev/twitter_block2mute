package main

import "github.com/faciam_dev/twitter_block2mute/backend/infrastructure"

func main() {
    db := infrastructure.NewDB()
    r := infrastructure.NewRouting(db)
    r.Run()
}