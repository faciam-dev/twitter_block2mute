package twitterapi

import (
	"context"
	"log"

	"github.com/ChimeraCoder/anaconda"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/user/userlookup"
	"github.com/michimani/gotwi/user/userlookup/types"
)

type GotwiHandler struct {
	ConsumerKey    string
	ConsumerSecret string
	CallbackUrl    string
	Api            *gotwi.Client
	AuthApi        *anaconda.TwitterApi
}

func NewGotwiHandler(config *config.Config) handler.TwitterHandler {
	return newGotwiHandler(GotwiHandler{
		ConsumerKey:    config.Twitter.ConsumerKey,
		ConsumerSecret: config.Twitter.ConsumerSecret,
		CallbackUrl:    config.Twitter.CallbackUrl,
	})
}

func newGotwiHandler(g GotwiHandler) handler.TwitterHandler {
	gotwiHandler := new(GotwiHandler)

	// for anaconda
	anaconda.SetConsumerKey(g.ConsumerKey)
	anaconda.SetConsumerSecret(g.ConsumerSecret)
	gotwiHandler.AuthApi = anaconda.NewTwitterApi("", "")

	// for gotwi
	in := &gotwi.NewClientInput{}

	c, err := gotwi.NewClient(in)
	if err != nil {
		log.Println(err)
		return gotwiHandler
	}
	gotwiHandler.Api = c

	return gotwiHandler
}

func (g *GotwiHandler) UpdateTwitterApi(token string, secret string) {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           token,
		OAuthTokenSecret:     secret,
	}

	c, err := gotwi.NewClient(in)
	if err != nil {
		log.Println(err)
		return
	}

	g.Api = c
}

func (g *GotwiHandler) AuthorizationURL() (string, error) {
	uri, _, err := g.AuthApi.AuthorizationURL(g.CallbackUrl)

	return uri, err
}

func (g *GotwiHandler) GetUser(twitterID string) (handler.TwitterUser, error) {

	p := &types.GetInput{ID: twitterID}

	user, err := userlookup.Get(context.Background(), g.Api, p)

	return TwitterUser{
		ID:          *user.Data.ID,
		ScreenName:  *user.Data.Name,
		TwitterName: *user.Data.Username,
	}, err
}

func (g *GotwiHandler) GetCredentials(token string, secret string) (handler.TwitterCredentials, handler.TwitterValues, error) {
	twitterCredentials := TwitterCredentials{}
	twitterValues := TwitterValues{}
	credentials, values, err := g.AuthApi.GetCredentials(&oauth.Credentials{
		Token: token,
	}, secret)

	twitterCredentials.Credentials = credentials
	twitterValues.Vaules = values

	if err != nil {
		return twitterCredentials, twitterValues, err
	}

	return twitterCredentials, twitterValues, nil
}
