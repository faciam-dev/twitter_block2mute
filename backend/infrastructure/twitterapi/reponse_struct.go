package twitterapi

import (
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
)

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

// TwitterUser
type TwitterUser struct {
	ID          string
	ScreenName  string
	TwitterName string
}

func (t TwitterUser) GetTwitterID() string {
	return t.ID
}

func (t TwitterUser) GetTwitterScreenName() string {
	return t.ScreenName
}

func (t TwitterUser) GetTwitterName() string {
	return t.TwitterName
}
