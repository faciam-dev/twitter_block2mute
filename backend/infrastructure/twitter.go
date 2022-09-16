package infrastructure

import "github.com/ChimeraCoder/anaconda"

type Twitter struct {
	ConsumerKey string
	ConsumerSecret string
	CallbackUrl string
	Api *anaconda.TwitterApi
}

func NewTwitter() *Twitter {
    c := NewConfig()
    return newTwitter(&Twitter {
		ConsumerKey: c.Twitter.Production.ConsumerKey,
		ConsumerSecret: c.Twitter.Production.ConsumerSecret,
		CallbackUrl: c.Twitter.Production.CallbackUrl,
	})
}

func newTwitter(t *Twitter) *Twitter {
	anaconda.SetConsumerKey(t.ConsumerKey)
	anaconda.SetConsumerSecret(t.ConsumerSecret)
	 
	t.Api = anaconda.NewTwitterApi("", "")

	return t
}