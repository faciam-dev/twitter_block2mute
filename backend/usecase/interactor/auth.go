package interactor

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
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

// auth url取得のusecase
func (a *Auth) Auth() {
	url, err := a.AuthRepo.GetAuthUrl()
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}
	auth := entity.NewAuth(url)
	a.OutputPort.RenderAuth(auth)
}

// auth判定のusecase
func (a *Auth) IsAuthenticated() {
	if err := a.AuthRepo.IsAuthenticated(); err != nil {
		a.LoggerHandler.Errorw("IsAuthenticated() error", "error", err.Error())
		a.OutputPort.RenderError(err)
		return
	}
	auth := entity.NewAuth("")
	auth.SuccessAuthenticated()
	a.OutputPort.RenderIsAuth(auth)
}

// callback処理のusecase
func (a *Auth) Callback(token string, secret string) {
	twitterCredentials, twitterValues, err := a.AuthRepo.AuthByCallbackParams(token, secret)
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}

	user, err := a.AuthRepo.FindUserByTwitterID(twitterValues.GetTwitterID())
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}

	// 新規ユーザー向け処理
	if user.GetID() == 0 {
		user.Update(
			user.GetID(),
			twitterValues.GetTwitterScreenName(),
			twitterValues.GetTwitterScreenName(),
			twitterValues.GetTwitterID(),
		)

		/*
			user.Name = twitterValues.GetTwitterScreenName()
			user.AccountName = twitterValues.GetTwitterScreenName()
			user.TwitterID = twitterValues.GetTwitterID()
		*/
	}

	if err := a.AuthRepo.UpsertUser(user); err != nil {
		a.OutputPort.RenderError(err)
		return
	}

	a.AuthRepo.UpdateTwitterApi(twitterCredentials.GetToken(), twitterCredentials.GetSecret())

	if err := a.AuthRepo.UpdateSession(twitterCredentials.GetToken(), twitterCredentials.GetSecret(), int(user.GetID()), user.GetTwitterID()); err != nil {
		a.OutputPort.RenderError(err)
		return
	}

	auth := entity.NewAuth("")
	auth.SuccessAuthenticated()

	a.OutputPort.RenderCallback(auth)
}

// logout処理のusecase
func (a *Auth) Logout() {
	err := a.AuthRepo.Logout()
	if err != nil {
		a.OutputPort.RenderError(err)
		return
	}
	auth := entity.NewAuth("")
	auth.SuccessLogout()
	a.OutputPort.RenderLogout(auth)
}
