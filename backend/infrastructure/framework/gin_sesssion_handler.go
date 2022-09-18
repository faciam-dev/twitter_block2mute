package framework

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type GinSessionHandler struct {
	Secret string
	Name string
	Gin *gin.Engine
	context *gin.Context
}

func NewGinSessionHandler(config *config.Config, gin *gin.Engine) gateway.SessionHandler {
    return newSessionHandler(GinSessionHandler{
		Secret: config.Session.Secret,
		Name: config.Session.Name,
		Gin: gin,
    })
}

func newSessionHandler(sessionHandler GinSessionHandler) gateway.SessionHandler {
    store := cookie.NewStore([]byte(sessionHandler.Secret))
    sessionHandler.Gin.Use(sessions.Sessions(sessionHandler.Name, store))
	return &sessionHandler
}

func (g *GinSessionHandler) SetContext(context interface{}) {
	g.context = context.(*gin.Context)
}

func (g *GinSessionHandler) Set(key string, value string) {
	session := g.getSession()
	session.Set(key, value)
}

func (g *GinSessionHandler) Get(key string) interface{} {
	session := g.getSession()
	value := session.Get(key)
	if value == nil {
		return nil
	}
	return value
}

func (g *GinSessionHandler) Save() error{
	session := g.getSession()
	return session.Save()
}

func (g *GinSessionHandler) getSession() sessions.Session {
	session := sessions.Default(g.context)
	return session
}