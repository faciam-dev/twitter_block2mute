package presenter

import (
	"net/http"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Auth struct {
	contextHandler gateway.ContextHandler
}

// NewUserOutputPort はUserOutputPortを取得します．
func NewAuthOutputPort(contextHandler gateway.ContextHandler) port.AuthOutputPort {
	return &Auth{
		contextHandler: contextHandler,
	}
}

// usecase.UserOutputPortを実装している
// Render はUserモデルを出力します．
func (a *Auth) RenderAuth(auth *entity.Auth) {
	a.contextHandler.JSON(http.StatusOK, map[string]interface{}{
		"to_url" : auth.AuthUrl, 
	})
}

func (a *Auth) RenderIsAuth(auth *entity.Auth) {
	a.contextHandler.JSON(http.StatusOK, map[string]interface{}{
		"result" : auth.Authenticated,
	})
}

// 空要素を返す
func (a *Auth) RenderCallback(auth *entity.Auth) {
	a.contextHandler.JSON(http.StatusOK, map[string]interface{}{
	})
}

// RenderError はErrorを出力します．
func (a *Auth) RenderError(err error) {
	a.contextHandler.JSON(http.StatusInternalServerError, map[string]interface{}{
		"error" : err,
	})
}