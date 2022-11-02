package controller

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Auth struct {
	// -> presenter.NewAuthOutputPort
	OutputFactory func(
		contextHandler handler.ContextHandler,
		LoggerHandler handler.LoggerHandler,
	) port.AuthOutputPort
	// -> interactor.NewAuthInputPort
	InputFactory func(
		outputPort port.AuthOutputPort,
		repository port.AuthRepository,
		logger handler.LoggerHandler,
	) port.AuthInputPort
	// -> gateway.NewAuthRepository
	RepoFactory func(
		contextHandler handler.ContextHandler,
		LoggerHandler handler.LoggerHandler,
		twitterHandler handler.TwitterHandler,
		sessionHandler handler.SessionHandler,
		UserDbHandler handler.UserDbHandler,
	) port.AuthRepository
	LoggerHandler  handler.LoggerHandler
	TwitterHandler handler.TwitterHandler
	SessionHandler handler.SessionHandler
	ContextHandler handler.ContextHandler
	UserDbHandler  handler.UserDbHandler
}

// 認証済みかどうかをセッションと照らし合わせて判別する
func (a *Auth) IsAuth(contextHandler handler.ContextHandler) {
	outputPort := a.OutputFactory(contextHandler, a.LoggerHandler)
	repository := a.RepoFactory(contextHandler, a.LoggerHandler, a.TwitterHandler, a.SessionHandler, a.UserDbHandler)
	inputPort := a.InputFactory(outputPort, repository, a.LoggerHandler)
	inputPort.IsAuthenticated()
}

// 認証を実施する
func (a *Auth) Auth(contextHandler handler.ContextHandler) {
	outputPort := a.OutputFactory(contextHandler, a.LoggerHandler)
	repository := a.RepoFactory(contextHandler, a.LoggerHandler, a.TwitterHandler, a.SessionHandler, a.UserDbHandler)
	inputPort := a.InputFactory(outputPort, repository, a.LoggerHandler)
	inputPort.Auth()
}

// コールバック処理
func (a *Auth) Callback(contextHandler handler.ContextHandler) {
	token := contextHandler.Query("oauth_token")
	secret := contextHandler.Query("oauth_verifier")
	twitterID := contextHandler.Query("user_id")
	twitterName := contextHandler.Query("screen_name")

	outputPort := a.OutputFactory(contextHandler, a.LoggerHandler)
	repository := a.RepoFactory(contextHandler, a.LoggerHandler, a.TwitterHandler, a.SessionHandler, a.UserDbHandler)
	inputPort := a.InputFactory(outputPort, repository, a.LoggerHandler)
	inputPort.Callback(token, secret, twitterID, twitterName)
}

// ログアウト処理
func (a *Auth) Logout(contextHandler handler.ContextHandler) {
	outputPort := a.OutputFactory(contextHandler, a.LoggerHandler)
	repository := a.RepoFactory(contextHandler, a.LoggerHandler, a.TwitterHandler, a.SessionHandler, a.UserDbHandler)
	inputPort := a.InputFactory(outputPort, repository, a.LoggerHandler)
	inputPort.Logout()
}
