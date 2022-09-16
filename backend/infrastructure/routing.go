package infrastructure

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/controller"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/presenter"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/interactor"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Routing struct {
    DB *DB
    Gin *gin.Engine
    Port string
    Twitter *Twitter
}

func NewRouting(db *DB, twitter *Twitter) *Routing {
    c := NewConfig()
    r := &Routing{
        DB: db,
        Gin: gin.Default(),
        Port: c.Routing.Port,
        Twitter: twitter,
    }
    r.setRouting()
    r.initSession(c.Session.Secret, c.Session.Name)
    return r
}

func (r *Routing) setRouting() {
	userController := controller.User{
		OutputFactory: presenter.NewUserOutputPort,
		InputFactory:  interactor.NewUserInputPort,
		RepoFactory:   gateway.NewUserRepository,
		Conn:          r.DB.Connection,
	}

    authController := controller.Auth{
		OutputFactory: presenter.NewAuthOutputPort,
		InputFactory:  interactor.NewAuthInputPort,
		RepoFactory:   gateway.NewAuthRepository,
		Conn:          r.DB.Connection,
        Api:           r.Twitter.Api,
        CallbackUrl:   r.Twitter.CallbackUrl,
    }

    // ルーティング割当
    r.Gin.GET("/user/user/:id", func (c *gin.Context) {userController.GetUserByID(c) })
    r.Gin.GET("/auth/auth", func (c *gin.Context) {authController.Auth(c) })
    r.Gin.GET("/auth/is_auth", func (c *gin.Context) {authController.IsAuth(c) })
    r.Gin.GET("/auth/auth_callback", func (c *gin.Context) {authController.Callback(c) })

}

func (r *Routing) initSession(sessionSecret string, sessionName string) {
    store := cookie.NewStore([]byte(sessionSecret))
    r.Gin.Use(sessions.Sessions(sessionName, store))
}

func (r *Routing) Run() {
    r.Gin.Run(r.Port)
}