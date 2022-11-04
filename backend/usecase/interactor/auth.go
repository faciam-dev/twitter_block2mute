package interactor

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Auth struct {
	OutputPort    port.AuthOutputPort
	AuthRepo      port.AuthRepository
	LoggerHandler handler.LoggerHandler
}

func NewAuthInputPort(
	outputPort port.AuthOutputPort,
	authRepository port.AuthRepository,
	loggerHandler handler.LoggerHandler,
) port.AuthInputPort {
	return &Auth{
		OutputPort:    outputPort,
		AuthRepo:      authRepository,
		LoggerHandler: loggerHandler,
	}
}

func (a *Auth) Auth() {
	auth, err := a.AuthRepo.Auth()
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}
	a.OutputPort.RenderAuth(auth)
}

func (a *Auth) IsAuthenticated() {
	auth, err := a.AuthRepo.IsAuthenticated()
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}
	a.OutputPort.RenderIsAuth(auth)
}

func (a *Auth) Callback(token string, secret string, twitterID string, twitterName string) {
	auth, err := a.AuthRepo.Callback(token, secret, twitterID, twitterName)
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}
	a.OutputPort.RenderCallback(auth)
}

func (a *Auth) Logout() {
	auth, err := a.AuthRepo.Logout()
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}
	a.OutputPort.RenderLogout(auth)
}
