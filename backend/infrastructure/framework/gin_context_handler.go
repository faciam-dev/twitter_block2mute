package framework

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/gin-gonic/gin"
)

type GinContextHandler struct {
    Context *gin.Context
}

func NewGinContextHandler(context *gin.Context) gateway.ContextHandler {
    ginContextHandler := GinContextHandler{
        Context: context,
    }

    return &ginContextHandler
}

func (g *GinContextHandler) GetContext() interface{} {
    return g.Context
}

func (g *GinContextHandler) Query(key string) (value string) {
    return g.Context.Query(key)
}

func (g *GinContextHandler) Param(key string) string {
    return g.Context.Param(key)
}

func (g *GinContextHandler) JSON(code int, obj interface{}) {
    g.Context.JSON(code, obj)
}
