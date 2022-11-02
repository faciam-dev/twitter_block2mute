package twitterapi

import (
	"context"
	"log"

	"github.com/ChimeraCoder/anaconda"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/user/block"
	blocktypes "github.com/michimani/gotwi/user/block/types"
	"github.com/michimani/gotwi/user/mute"
	mutetypes "github.com/michimani/gotwi/user/mute/types"
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

func (g *GotwiHandler) GetBlockedUser(twitterID string) (handler.TwitterUserIds, error) {

	tokenDefault := "NEXT"
	tokenEnd := "0"
	nextToken := tokenDefault
	total := 0
	twitterIDs := []string{}
	var err error

	for nextToken != tokenEnd {
		p := &blocktypes.ListInput{ID: twitterID}

		if nextToken != tokenDefault {
			p.PaginationToken = nextToken
		}

		listOutput, err := block.List(context.Background(), g.Api, p)

		if err != nil {
			return TwitterUserIds{
				IDs:   twitterIDs,
				Total: total,
			}, err
		}

		if listOutput.Meta.NextToken == nil {
			nextToken = tokenEnd
		} else {
			nextToken = *listOutput.Meta.NextToken
		}

		total += *listOutput.Meta.ResultCount

		for _, user := range listOutput.Data {
			twitterIDs = append(twitterIDs, *user.ID)
		}
	}

	return TwitterUserIds{
		IDs:   twitterIDs,
		Total: total,
	}, err
}

func (g *GotwiHandler) CreateMute(twitterID string, targetTwitterID string) error {

	p := &mutetypes.CreateInput{
		ID:       twitterID,
		TargetID: targetTwitterID,
	}

	_, err := mute.Create(context.Background(), g.Api, p)

	return err
}

func (g *GotwiHandler) DestroyBlock(twitterID string, targetTwitterID string) error {
	p := &blocktypes.DeleteInput{
		SourceUserID: twitterID,
		TargetID:     targetTwitterID,
	}

	_, err := block.Delete(context.Background(), g.Api, p)

	return err
}
