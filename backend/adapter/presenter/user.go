package presenter

import (
	"net/http"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
	"github.com/gin-gonic/gin"
)

type User struct {
	c *gin.Context
}

// NewUserOutputPort はUserOutputPortを取得します．
func NewUserOutputPort(c *gin.Context) port.UserOutputPort {
	return &User{
		c: c,
	}
}

// usecase.UserOutputPortを実装している
// Render はUserモデルを出力します．
func (u *User) Render(user *entity.User) {
	u.c.JSON(http.StatusOK, gin.H{
		"id" : user.ID,
		"name" : user.Name,
	})
}

// RenderNotFound は ユーザーがないことを出力します。
func (u *User) RenderNotFound() {
	u.c.JSON(http.StatusNotFound, gin.H{
		
	})
}

// RenderError はErrorを出力します．
func (u *User) RenderError(err error) {
	u.c.JSON(http.StatusInternalServerError, gin.H{
		"error" : err,
	})
}