package presenter

import (
	"net/http"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	ctx *gin.Context
}

// NewUserOutputPort はUserOutputPortを取得します．
func NewAuthOutputPort(c *gin.Context) port.AuthOutputPort {
	return &Auth{
		ctx: c,
	}
}

// usecase.UserOutputPortを実装している
// Render はUserモデルを出力します．
func (a *Auth) RenderAuth(auth *entity.Auth) {
	a.ctx.JSON(http.StatusOK, gin.H{
		"to_url" : auth.AuthUrl, 
	})
}

func (a *Auth) RenderIsAuth(auth *entity.Auth) {
	a.ctx.JSON(http.StatusOK, gin.H{
		"result" : auth.Authenticated,
	})
}

// 空要素を返す
func (a *Auth) RenderCallback(auth *entity.Auth) {
	a.ctx.JSON(http.StatusOK, gin.H{
	})
}

// RenderError はErrorを出力します．
func (a *Auth) RenderError(err error) {
	a.ctx.JSON(http.StatusInternalServerError, gin.H{
		"error" : err,
	})
}