package presenter

import (
	"net/http"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
	"github.com/gin-gonic/gin"
)

type User struct {
	ctx *gin.Context
}

// NewUserOutputPort はUserOutputPortを取得します．
func NewUserOutputPort(ctx *gin.Context) port.UserOutputPort {
	return &User{
		ctx: ctx,
	}
}

// usecase.UserOutputPortを実装している
// Render はUserモデルを出力します．
func (u *User) Render(user *entity.User) {
	u.ctx.JSON(http.StatusOK, gin.H{
		"id" : user.ID,
		"name" : user.Name,
	})
}

// RenderNotFound は ユーザーがないことを出力します。
func (u *User) RenderNotFound() {
	u.ctx.JSON(http.StatusNotFound, gin.H{
		
	})
}

// RenderError はErrorを出力します．
func (u *User) RenderError(err error) {
	u.ctx.JSON(http.StatusInternalServerError, gin.H{
		"error" : err,
	})
}