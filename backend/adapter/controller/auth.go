package controller

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Auth struct {
	// -> presenter.NewAuthOutputPort
	OutputFactory func(contextHandler gateway.ContextHandler) port.AuthOutputPort
	// -> interactor.NewAuthInputPort
	InputFactory func(o port.AuthOutputPort, u port.AuthRepository) port.AuthInputPort
	// -> gateway.NewAuthRepository
	RepoFactory func(contextHandler gateway.ContextHandler, twitterHandler gateway.TwitterHandler, sessionHandler gateway.SessionHandler) port.AuthRepository
	TwitterHandler gateway.TwitterHandler
	SessionHandler gateway.SessionHandler
	ContextHandler gateway.ContextHandler
}

// 認証済みかどうかをセッションと照らし合わせて判別する
func (a *Auth) IsAuth(contextHandler gateway.ContextHandler) {
	outputPort := a.OutputFactory(contextHandler)
	repository := a.RepoFactory(contextHandler, a.TwitterHandler, a.SessionHandler)
	inputPort := a.InputFactory(outputPort, repository)
	inputPort.IsAuthenticated()
}

// 認証を実施する
func (a *Auth) Auth(contextHandler gateway.ContextHandler) {
	outputPort := a.OutputFactory(contextHandler)
	repository := a.RepoFactory(contextHandler, a.TwitterHandler, a.SessionHandler)
	inputPort := a.InputFactory(outputPort, repository)
	inputPort.Auth()
}

// コールバック処理
func (a *Auth) Callback(contextHandler gateway.ContextHandler) {
	token := contextHandler.Query("oauth_token")
    secret := contextHandler.Query("oauth_verifier")

	outputPort := a.OutputFactory(contextHandler)
	repository := a.RepoFactory(contextHandler, a.TwitterHandler, a.SessionHandler)
	inputPort := a.InputFactory(outputPort, repository)
	inputPort.Callback(token, secret)
}