package controller

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	// -> presenter.NewAuthOutputPort
	OutputFactory func(c *gin.Context) port.AuthOutputPort
	// -> interactor.NewAuthInputPort
	InputFactory func(o port.AuthOutputPort, u port.AuthRepository) port.AuthInputPort
	// -> gateway.NewAuthRepository
	RepoFactory func(ctx *gin.Context, twitterHandler gateway.TwitterHandler, sessionHandler gateway.SessionHandler) port.AuthRepository
	TwitterHandler gateway.TwitterHandler
	SessionHandler gateway.SessionHandler
	//Conn *gorm.DB
	//Api *anaconda.TwitterApi
	//CallbackUrl string
}

// 認証済みかどうかをセッションと照らし合わせて判別する
func (a *Auth) IsAuth(c *gin.Context) {
	outputPort := a.OutputFactory(c)
	repository := a.RepoFactory(c, a.TwitterHandler, a.SessionHandler)
	inputPort := a.InputFactory(outputPort, repository)
	inputPort.IsAuthenticated()
}

// 認証を実施する
func (a *Auth) Auth(c *gin.Context) {
	outputPort := a.OutputFactory(c)
	repository := a.RepoFactory(c, a.TwitterHandler, a.SessionHandler)
	inputPort := a.InputFactory(outputPort, repository)
	inputPort.Auth()
}

// コールバック処理
func (a *Auth) Callback(c *gin.Context) {
	outputPort := a.OutputFactory(c)
	repository := a.RepoFactory(c, a.TwitterHandler, a.SessionHandler)
	inputPort := a.InputFactory(outputPort, repository)
	inputPort.Callback()
}