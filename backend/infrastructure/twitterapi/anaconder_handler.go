package twitterapi

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/garyburd/go-oauth/oauth"
)

type AnaconderHandler struct {
	ConsumerKey string
	ConsumerSecret string
	CallbackUrl string
	Api *anaconda.TwitterApi
}

func NewAnaconderHandler(config *config.Config) handler.TwitterHandler {
    return newAnaconderHandler(AnaconderHandler{
		ConsumerKey: config.Twitter.ConsumerKey,
		ConsumerSecret: config.Twitter.ConsumerSecret,
		CallbackUrl: config.Twitter.CallbackUrl,
    })
}

func newAnaconderHandler(a AnaconderHandler) handler.TwitterHandler {
	anaconda.SetConsumerKey(a.ConsumerKey)
	anaconda.SetConsumerSecret(a.ConsumerSecret)
	
	anaconderHandler := new(AnaconderHandler)

	anaconderHandler.Api = anaconda.NewTwitterApi("", "")

	return anaconderHandler
}

func (a *AnaconderHandler) SetCredentials(token string, secret string) {
	a.Api.Credentials.Token = token
	a.Api.Credentials.Secret = secret
}

func (a *AnaconderHandler) AuthorizationURL() (string, error) {
	uri, _, err := a.Api.AuthorizationURL(a.CallbackUrl)

	return uri, err
}

func (a *AnaconderHandler) GetRateLimits() error {
	r := []string{}
	_, err := a.Api.GetRateLimits(r)

	return err
}

func (a *AnaconderHandler) GetCredentials(token string, secret string) (handler.TwitterCredentials, error) {
	twitterCredentials := TwitterCredentials{}
	credentials, _, err := a.Api.GetCredentials(&oauth.Credentials{
        Token: token,
    }, secret)

	twitterCredentials.Credentials = credentials

	if err != nil {
        return twitterCredentials, err
    }

	return twitterCredentials, nil
}

type TwitterCredentials struct {
	Credentials *oauth.Credentials
}

func (t TwitterCredentials) GetToken() string {
	return t.Credentials.Token
}

func (t TwitterCredentials) GetSecret() string {
	return t.Credentials.Secret
}