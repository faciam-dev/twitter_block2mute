package port

import (
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

type AuthInputPort interface {
	IsAuthenticated(consumerKey string, consumerSecret string)
	Auth(consumerKey string, consumerSecret string, callbackUrl string)
}

type AuthOutputPort interface {
	Render(*entity.Auth)
	RenderError(error)
}

type AuthRepository interface {
	IsAuthenticated(consumerKey string, consumerSecret string) (*entity.Auth, error)
	Auth(consumerKey string, consumerSecret string, callbackUrl string) (*entity.Auth, error)
	Callback(consumerKey string, consumerSecret string) (*entity.Auth, error)
}