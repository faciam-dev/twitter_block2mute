package presenter

import (
	"net/http"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type User struct {
	contextHandler handler.ContextHandler
	loggerHandler  handler.LoggerHandler
}

// NewUserOutputPort はUserOutputPortを取得します．
func NewUserOutputPort(contextHandler handler.ContextHandler, loggerHandler handler.LoggerHandler) port.UserOutputPort {
	return &User{
		contextHandler: contextHandler,
		loggerHandler:  loggerHandler,
	}
}

// usecase.UserOutputPortを実装している
// Render はUserモデルを出力します．
func (u *User) Render(user *entity.User) {
	u.contextHandler.JSON(http.StatusOK, map[string]interface{}{
		"id":   user.GetID(),
		"name": user.GetName(),
	})
}

// RenderNotFound は ユーザーがないことを出力します。
func (u *User) RenderNotFound() {
	u.contextHandler.JSON(http.StatusNotFound, map[string]interface{}{})
}

// RenderError はErrorを出力します．
func (u *User) RenderError(err error) {
	u.contextHandler.JSON(http.StatusInternalServerError, map[string]interface{}{
		"error": err,
	})
}
