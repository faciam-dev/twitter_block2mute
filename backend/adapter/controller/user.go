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
	RepoFactory func(dbHandler handler.DbHandler) port.UserRepository
	DbHandler handler.DbHandler
}

// GetUserByID は，httpを受け取り，portを組み立てて，inputPort.GetUserByIDを呼び出します．
func (u *User) GetUserByID(contextHandler handler.ContextHandler) {

	//id, _ := strconv.Atoi(c.Param("id"))
	id := contextHandler.Param("id")

	outputPort := u.OutputFactory(contextHandler)
	repository := u.RepoFactory(u.DbHandler)
	inputPort := u.InputFactory(outputPort, repository)
	inputPort.GetUserByID(id)

}