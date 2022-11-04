package controller

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type User struct {
	// -> presenter.NewUserOutputPort
	OutputFactory func(
		contextHandler handler.ContextHandler,
		LoggerHandler handler.LoggerHandler,
	) port.UserOutputPort

	// -> interactor.NewUserInputPort
	InputFactory func(
		o port.UserOutputPort,
		u port.UserRepository,
		LoggerHandler handler.LoggerHandler,
	) port.UserInputPort

	// -> gateway.NewUserRepository
	RepoFactory func(
		LoggerHandler handler.LoggerHandler,
		dbHandler handler.UserDbHandler,
	) port.UserRepository

	LoggerHandler  handler.LoggerHandler
	SessionHandler handler.SessionHandler
	UserDbHandler  handler.UserDbHandler
}

// GetUserByID は，httpを受け取り，portを組み立てて，inputPort.GetUserByIDを呼び出します．
func (u *User) GetUserSelf(contextHandler handler.ContextHandler) {
	u.SessionHandler.SetContextHandler(contextHandler)
	userID := u.SessionHandler.Get("user_id")

	// nilの場合は空のIDをあたえ、エラーにさせる
	if userID == nil {
		userID = ""
	}

	outputPort := u.OutputFactory(contextHandler, u.LoggerHandler)
	repository := u.RepoFactory(u.LoggerHandler, u.UserDbHandler)
	inputPort := u.InputFactory(outputPort, repository, u.LoggerHandler)
	inputPort.GetUserByID(userID.(string))

}
