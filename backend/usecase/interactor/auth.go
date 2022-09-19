package interactor

import (
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Auth struct {
	OutputPort port.AuthOutputPort
	AuthRepo   port.AuthRepository
}

func NewAuthInputPort(outputPort port.AuthOutputPort, authRepository port.AuthRepository) port.AuthInputPort {
	return &Auth{
		OutputPort: outputPort,
		AuthRepo:   authRepository,
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
	auth, err := a.AuthRepo.IsAuthenticated();
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}
	a.OutputPort.RenderIsAuth(auth)
}

func (a *Auth) Callback(token string, secret string, twitterID string, twitterName string) {
	auth, err := a.AuthRepo.Callback(token, secret, twitterID, twitterName);
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}
	a.OutputPort.RenderCallback(auth)
}