package port

import (
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

type AuthInputPort interface {
	IsAuthenticated()
	Auth()
	Callback(token string, secret string, twitterID string, twitterName string)
}

type AuthOutputPort interface {
	RenderAuth(*entity.Auth)
	RenderIsAuth(*entity.Auth)
	RenderCallback(*entity.Auth)
	RenderError(error)
}

type AuthRepository interface {
	IsAuthenticated() (*entity.Auth, error)
	Auth() (*entity.Auth, error)
	Callback(token string, secret string, twitterID string, twitterName string) (*entity.Auth, error)
}