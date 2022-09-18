package framework

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/controller"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/presenter"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/interactor"
	"github.com/gin-gonic/gin"
)

type Routing struct {
    config *config.Config
    Gin *gin.Engine
    Port string
}

func NewRouting(config *config.Config, dbHandler handler.DbHandler, twitterHandler handler.TwitterHandler) *Routing {
    r := &Routing{
        config: config,
        Gin: gin.Default(),
        Port: config.Routing.Port,
    }
    sessionHandler := NewGinSessionHandler(config, r.Gin)
    r.setRouting(dbHandler, twitterHandler, sessionHandler)
    return r
}

func (r *Routing) setRouting(dbHandler handler.DbHandler, twitterHandler handler.TwitterHandler, sessionHandler handler.SessionHandler) {
	userController := controller.User{
		OutputFactory: presenter.NewUserOutputPort,
		InputFactory:  interactor.NewUserInputPort,
		RepoFactory:   gateway.NewUserRepository,
        DbHandler:     dbHandler,
	}

    authController := controller.Auth{
		OutputFactory:  presenter.NewAuthOutputPort,
		InputFactory:   interactor.NewAuthInputPort,
		RepoFactory:    gateway.NewAuthRepository,
        TwitterHandler: twitterHandler,
        SessionHandler: sessionHandler,
    }

    // ルーティング割当
    // user
    r.Gin.GET("/user/user/:id", func (c *gin.Context) {
        userController.GetUserByID(NewGinContextHandler(c))
    })
    // auth
    r.Gin.GET("/auth/auth", func (c *gin.Context) {
        authController.Auth(NewGinContextHandler(c))
    })
    r.Gin.GET("/auth/is_auth", func (c *gin.Context) {
        authController.IsAuth(NewGinContextHandler(c))
    })
    r.Gin.GET("/auth/auth_callback", func (c *gin.Context) {
        authController.Callback(NewGinContextHandler(c))
    })
    // block2mute

}

func (r *Routing) Run() {
    r.Gin.Run(r.Port)
}