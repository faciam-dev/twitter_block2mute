package presenter

import (
	"net/http"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type User struct {
	contextHandler gateway.ContextHandler
}

// NewUserOutputPort はUserOutputPortを取得します．
func NewUserOutputPort(contextHandler gateway.ContextHandler) port.UserOutputPort {
	return &User{
		contextHandler: contextHandler,
	}
}

// usecase.UserOutputPortを実装している
// Render はUserモデルを出力します．
func (u *User) Render(user *entity.User) {
	u.contextHandler.JSON(http.StatusOK, map[string]interface{}{
		"id" : user.ID,
		"name" : user.Name,
	})
}

// RenderNotFound は ユーザーがないことを出力します。
func (u *User) RenderNotFound() {
	u.contextHandler.JSON(http.StatusNotFound, map[string]interface{}{
		
	})
}

// RenderError はErrorを出力します．
func (u *User) RenderError(err error) {
	u.contextHandler.JSON(http.StatusInternalServerError, map[string]interface{}{
		"error" : err,
	})
}