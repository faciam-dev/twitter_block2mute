package controller

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Auth struct {
	// -> presenter.NewAuthOutputPort
	OutputFactory func(contextHandler handler.ContextHandler) port.AuthOutputPort
	// -> interactor.NewAuthInputPort
	InputFactory func(o port.AuthOutputPort, u port.AuthRepository) port.AuthInputPort
	// -> gateway.NewAuthRepository
	RepoFactory func(contextHandler handler.ContextHandler, twitterHandler handler.TwitterHandler, sessionHandler handler.SessionHandler) port.AuthRepository
	TwitterHandler handler.TwitterHandler
	SessionHandler handler.SessionHandler
	ContextHandler handler.ContextHandler
}

// 認証済みかどうかをセッションと照らし合わせて判別する
func (a *Auth) IsAuth(contextHandler handler.ContextHandler) {
	outputPort := a.OutputFactory(contextHandler)
	repository := a.RepoFactory(contextHandler, a.TwitterHandler, a.SessionHandler)
	inputPort := a.InputFactory(outputPort, repository)
	inputPort.IsAuthenticated()
}

// 認証を実施する
func (a *Auth) Auth(contextHandler handler.ContextHandler) {
	outputPort := a.OutputFactory(contextHandler)
	repository := a.RepoFactory(contextHandler, a.TwitterHandler, a.SessionHandler)
	inputPort := a.InputFactory(outputPort, repository)
	inputPort.Auth()
}

// コールバック処理
func (a *Auth) Callback(contextHandler handler.ContextHandler) {
	token := contextHandler.Query("oauth_token")
    secret := contextHandler.Query("oauth_verifier")

	outputPort := a.OutputFactory(contextHandler)
	repository := a.RepoFactory(contextHandler, a.TwitterHandler, a.SessionHandler)
	inputPort := a.InputFactory(outputPort, repository)
	inputPort.Callback(token, secret)
}