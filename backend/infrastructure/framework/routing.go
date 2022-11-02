package framework

import (
	"net/http"
	"strings"
	"time"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/controller"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/presenter"
	"github.com/faciam_dev/twitter_block2mute/backend/common"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/interactor"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

type Routing struct {
	config *config.Config
	Gin    *gin.Engine
	Port   string
}

func NewRouting(
	config *config.Config,
	loggerHandler handler.LoggerHandler,
	dbHandler database.GormDbHandler,
	twitterHandler handler.TwitterHandler,
) *Routing {
	r := &Routing{
		config: config,
		Gin:    gin.Default(),
		Port:   config.Routing.Port,
	}
	r.GinModeConfig()
	r.AllowOrigins() // before set routing
	r.CsrfConfig()
	sessionHandler := NewGinSessionHandler(config, r.Gin)
	r.setRouting(loggerHandler, dbHandler, twitterHandler, sessionHandler)
	return r
}

func (r *Routing) setRouting(
	loggerHandler handler.LoggerHandler,
	dbHandler database.GormDbHandler,
	twitterHandler handler.TwitterHandler,
	sessionHandler handler.SessionHandler,
) {
	userController := controller.User{
		OutputFactory:  presenter.NewUserOutputPort,
		InputFactory:   interactor.NewUserInputPort,
		RepoFactory:    gateway.NewUserRepository,
		SessionHandler: sessionHandler,
		UserDbHandler:  database.NewUserDbHandler(dbHandler),
	}

	authController := controller.Auth{
		OutputFactory:  presenter.NewAuthOutputPort,
		InputFactory:   interactor.NewAuthInputPort,
		RepoFactory:    gateway.NewAuthRepository,
		LoggerHandler:  loggerHandler,
		TwitterHandler: twitterHandler,
		SessionHandler: sessionHandler,
		UserDbHandler:  database.NewUserDbHandler(dbHandler),
	}

	blockController := controller.Block{
		OutputFactory:  presenter.NewBlockOutputPort,
		InputFactory:   interactor.NewBlockInputPort,
		RepoFactory:    gateway.NewBlockRepository,
		TwitterHandler: twitterHandler,
		SessionHandler: sessionHandler,
		BlockDbHandler: database.NewBlockDbHandler(dbHandler),
		UserDbHandler:  database.NewUserDbHandler(dbHandler),
	}

	block2MuteController := controller.Block2Mute{
		OutputFactory:  presenter.NewBlock2MuteOutputPort,
		InputFactory:   interactor.NewBlock2MuteInputPort,
		RepoFactory:    gateway.NewBlock2MuteRepository,
		TwitterHandler: twitterHandler,
		SessionHandler: sessionHandler,
		BlockDbHandler: database.NewBlockDbHandler(dbHandler),
		UserDbHandler:  database.NewUserDbHandler(dbHandler),
		MuteDbHandler:  database.NewMuteHandler(dbHandler),
	}

	// PROXY
	r.Gin.SetTrustedProxies(r.config.Routing.TrustedProxies)

	// ルーティング割当
	// session
	r.Gin.GET("/session", func(c *gin.Context) {
		c.Header("X-CSRF-Token", csrf.Token(c.Request))
		c.JSON(http.StatusOK, map[string]interface{}{})
	})
	// user
	r.Gin.POST("/user/user/self", func(c *gin.Context) {
		userController.GetUserSelf(NewGinContextHandler(c))
	})
	// auth
	r.Gin.POST("/auth/auth", func(c *gin.Context) {
		authController.Auth(NewGinContextHandler(c))
	})
	r.Gin.POST("/auth/is_auth", func(c *gin.Context) {
		authController.IsAuth(NewGinContextHandler(c))
	})
	r.Gin.POST("/auth/logout", func(c *gin.Context) {
		authController.Logout(NewGinContextHandler(c))
	})
	r.Gin.GET("/auth/auth_callback", func(c *gin.Context) {
		authController.Callback(NewGinContextHandler(c))
	})
	// block
	r.Gin.POST("/blocks/ids", func(c *gin.Context) {
		blockController.GetBlockByID(NewGinContextHandler(c))
	})
	// block2mute
	r.Gin.POST("/block2mute/all", func(c *gin.Context) {
		block2MuteController.All(NewGinContextHandler(c))
	})
}

func (r *Routing) AllowOrigins() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = r.config.Routing.AllowOrigins
	corsConfig.AllowHeaders = r.config.Routing.AllowHeaders
	corsConfig.ExposeHeaders = r.config.Routing.ExposeHeaders
	corsConfig.AllowMethods = r.config.Routing.AllowMethods
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = time.Duration(r.config.Routing.MaxAge)
	r.Gin.Use(cors.New(corsConfig))
}

func (r *Routing) CsrfConfig() {
	csrfMiddleware := csrf.Protect(
		[]byte(common.RandomString(32)),
		csrf.TrustedOrigins(r.config.Routing.AllowOrigins),
		csrf.Secure(r.config.Routing.CsrfSecure),
		csrf.Path("/"),
	)
	r.Gin.Use(adapter.Wrap(csrfMiddleware))
}

func (r *Routing) GinModeConfig() {
	if strings.ToLower(r.config.ReleaseMode) == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func (r *Routing) Run() {
	r.Gin.Run(r.Port)
}
