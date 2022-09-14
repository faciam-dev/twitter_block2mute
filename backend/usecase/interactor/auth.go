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

func (a *Auth) Auth(consumerKey string, consumerSecret string, callbackUrl string) {
	auth, err := a.AuthRepo.Auth(consumerKey, consumerSecret, callbackUrl)
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}
	a.OutputPort.Render(auth)
}

func (a *Auth) IsAuthenticated(consumerKey string, consumerSecret string) {
	auth, err := a.AuthRepo.IsAuthenticated(consumerKey, consumerSecret);
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}
	a.OutputPort.Render(auth)
}