package port

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

type AuthInputPort interface {
	IsAuthenticated()
	Auth()
	Callback(token string, secret string)
	Logout()
}

type AuthOutputPort interface {
	RenderAuth(*entity.Auth)
	RenderIsAuth(*entity.Auth)
	RenderCallback(*entity.Auth)
	RenderLogout(*entity.Auth)
	RenderError(error)
}

type AuthRepository interface {
	IsAuthenticated() error
	GetAuthUrl() (string, error)
	AuthByCallbackParams(token string, secret string) (handler.TwitterCredentials, handler.TwitterValues, error)
	FindUserByTwitterID(twitterID string) (*entity.User, error)
	UpsertUser(user *entity.User) error
	UpdateTwitterApi(token string, secret string)
	UpdateSession(token string, secret string, userID int, twitterID string) error
	Logout() error
}
