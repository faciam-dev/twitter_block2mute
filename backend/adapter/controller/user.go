package controller

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type User struct {
	// -> presenter.NewUserOutputPort
	OutputFactory func(contextHandler gateway.ContextHandler) port.UserOutputPort
	// -> interactor.NewUserInputPort
	InputFactory func(o port.UserOutputPort, u port.UserRepository) port.UserInputPort
	// -> gateway.NewUserRepository
	RepoFactory func(dbHandler gateway.DbHandler) port.UserRepository
	DbHandler gateway.DbHandler
}

// GetUserByID は，httpを受け取り，portを組み立てて，inputPort.GetUserByIDを呼び出します．
func (u *User) GetUserByID(contextHandler gateway.ContextHandler) {

	//id, _ := strconv.Atoi(c.Param("id"))
	id := contextHandler.Param("id")

	outputPort := u.OutputFactory(contextHandler)
	repository := u.RepoFactory(u.DbHandler)
	inputPort := u.InputFactory(outputPort, repository)
	inputPort.GetUserByID(id)

}