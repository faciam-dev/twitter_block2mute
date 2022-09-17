package infrastructure

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
)

type Twitter struct {
	ConsumerKey string
	ConsumerSecret string
	CallbackUrl string
	Api *anaconda.TwitterApi
}

func NewTwitter() *Twitter {
    c := config.NewConfig()
    return newTwitter(&Twitter {
		ConsumerKey: c.Twitter.ConsumerKey,
		ConsumerSecret: c.Twitter.ConsumerSecret,
		CallbackUrl: c.Twitter.CallbackUrl,
	})
}

func newTwitter(t *Twitter) *Twitter {
	anaconda.SetConsumerKey(t.ConsumerKey)
	anaconda.SetConsumerSecret(t.ConsumerSecret)
	 
	t.Api = anaconda.NewTwitterApi("", "")

	return t
}