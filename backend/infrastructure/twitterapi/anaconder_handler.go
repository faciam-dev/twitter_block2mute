package twitterapi

import (
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/garyburd/go-oauth/oauth"
)

type AnaconderHandler struct {
	ConsumerKey    string
	ConsumerSecret string
	CallbackUrl    string
	Api            *anaconda.TwitterApi
}

func NewAnaconderHandler(config *config.Config) handler.TwitterHandler {
	return newAnaconderHandler(AnaconderHandler{
		ConsumerKey:    config.Twitter.ConsumerKey,
		ConsumerSecret: config.Twitter.ConsumerSecret,
		CallbackUrl:    config.Twitter.CallbackUrl,
	})
}

func newAnaconderHandler(a AnaconderHandler) handler.TwitterHandler {
	anaconda.SetConsumerKey(a.ConsumerKey)
	anaconda.SetConsumerSecret(a.ConsumerSecret)

	anaconderHandler := new(AnaconderHandler)

	anaconderHandler.Api = anaconda.NewTwitterApi("", "")

	return anaconderHandler
}

func (a *AnaconderHandler) UpdateTwitterApi(token string, secret string) {
	a.Api = anaconda.NewTwitterApi(token, secret)
}

func (a *AnaconderHandler) AuthorizationURL() (string, error) {
	uri, _, err := a.Api.AuthorizationURL(a.CallbackUrl)

	return uri, err
}

func (a *AnaconderHandler) GetUser(twitterID string) (handler.TwitterUser, error) {
	values := url.Values{}
	convertedTwitterID, _ := strconv.ParseInt(twitterID, 10, 64)

	user, err := a.Api.GetUsersShowById(convertedTwitterID, values)

	return TwitterUser{
		ID:          user.IdStr,
		ScreenName:  user.ScreenName,
		TwitterName: user.Name,
	}, err
}

func (a *AnaconderHandler) GetCredentials(token string, secret string) (handler.TwitterCredentials, handler.TwitterValues, error) {
	twitterCredentials := TwitterCredentials{}
	twitterValues := TwitterValues{}
	credentials, values, err := a.Api.GetCredentials(&oauth.Credentials{
		Token: token,
	}, secret)

	twitterCredentials.Credentials = credentials
	twitterValues.Vaules = values

	if err != nil {
		return twitterCredentials, twitterValues, err
	}

	return twitterCredentials, twitterValues, nil
}

// TwitterCredentials
type TwitterCredentials struct {
	Credentials *oauth.Credentials
}

func (t TwitterCredentials) GetToken() string {
	return t.Credentials.Token
}

func (t TwitterCredentials) GetSecret() string {
	return t.Credentials.Secret
}

// TwitterValues
type TwitterValues struct {
	Vaules url.Values
}

func (t TwitterValues) GetTwitterID() string {
	return t.getTwitterValue("user_id")
}

func (t TwitterValues) GetTwitterScreenName() string {
	return t.getTwitterValue("screen_name")
}

func (t TwitterValues) getTwitterValue(key string) string {
	for k, v := range t.Vaules {
		if k == key {
			return v[0]
		}
	}
	return ""
}
