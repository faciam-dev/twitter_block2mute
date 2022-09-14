package controller

import (
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Twitter struct {
	ConsumerKey string
	ConsumerSecret string
	CallbackUrl string    
}

type Auth struct {
	// -> presenter.NewAuthOutputPort
	OutputFactory func(c *gin.Context) port.AuthOutputPort
	// -> interactor.NewAuthInputPort
	InputFactory func(o port.AuthOutputPort, u port.AuthRepository) port.AuthInputPort
	// -> gateway.NewAuthRepository
	RepoFactory func(c *gorm.DB) port.AuthRepository
	Conn *gorm.DB
	Twitter Twitter
}

// 認証済みかどうかをセッションと照らし合わせて判別する
func (a *Auth) isAuth(c *gin.Context) {
	outputPort := a.OutputFactory(c)
	repository := a.RepoFactory(a.Conn)
	inputPort := a.InputFactory(outputPort, repository)
	inputPort.IsAuthenticated(a.Twitter.ConsumerKey, a.Twitter.ConsumerSecret)
}

// 認証を実施する
func (a *Auth) Auth(c *gin.Context) {
	outputPort := a.OutputFactory(c)
	repository := a.RepoFactory(a.Conn)
	inputPort := a.InputFactory(outputPort, repository)
	inputPort.Auth(a.Twitter.ConsumerKey, a.Twitter.ConsumerSecret, a.Twitter.CallbackUrl)
}
