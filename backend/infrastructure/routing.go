package infrastructure

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/controller"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/presenter"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/framework"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/interactor"
	"github.com/gin-gonic/gin"
)

type Routing struct {
    config *config.Config
    Gin *gin.Engine
    Port string
}

func NewRouting(config *config.Config, dbHandler gateway.DbHandler, twitterHandler gateway.TwitterHandler) *Routing {
    r := &Routing{
        config: config,
        Gin: gin.Default(),
        Port: config.Routing.Port,
    }
    //r.initSession(config.Session.Secret, config.Session.Name)
    sessionHandler := framework.NewGinSessionHandler(config, r.Gin)
    r.setRouting(dbHandler, twitterHandler, sessionHandler)
    return r
}

func (r *Routing) setRouting(dbHandler gateway.DbHandler, twitterHandler gateway.TwitterHandler, sessionHandler gateway.SessionHandler) {
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
    r.Gin.GET("/user/user/:id", func (c *gin.Context) {userController.GetUserByID(c) })
    r.Gin.GET("/auth/auth", func (c *gin.Context) {authController.Auth(c) })
    r.Gin.GET("/auth/is_auth", func (c *gin.Context) {authController.IsAuth(c) })
    r.Gin.GET("/auth/auth_callback", func (c *gin.Context) {authController.Callback(c) })

}

func (r *Routing) Run() {
    r.Gin.Run(r.Port)
}