package controller

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type User struct {
	// -> presenter.NewUserOutputPort
	OutputFactory func(contextHandler handler.ContextHandler) port.UserOutputPort
	// -> interactor.NewUserInputPort
	InputFactory func(o port.UserOutputPort, u port.UserRepository) port.UserInputPort
	// -> gateway.NewUserRepository
	RepoFactory    func(dbHandler handler.UserDbHandler) port.UserRepository
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

	outputPort := u.OutputFactory(contextHandler)
	repository := u.RepoFactory(u.UserDbHandler)
	inputPort := u.InputFactory(outputPort, repository)
	inputPort.GetUserByID(userID.(string))

}
