package controller

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
	"github.com/gin-gonic/gin"
)

type User struct {
	// -> presenter.NewUserOutputPort
	OutputFactory func(ctx *gin.Context) port.UserOutputPort
	// -> interactor.NewUserInputPort
	InputFactory func(o port.UserOutputPort, u port.UserRepository) port.UserInputPort
	// -> gateway.NewUserRepository
	RepoFactory func(dbHandler gateway.DbHandler) port.UserRepository
	DbHandler gateway.DbHandler
}


/*
func NewUserController(dbHandler gateway.DbHandler) *User {
	userController := User{
		OutputFactory: presenter.NewUserOutputPort,
		InputFactory:  interactor.NewUserInputPort,
		RepoFactory:   gateway.NewUserRepository(dbHandler),
	}

    return &UserController{
        Interactor: usecase.UserInteractor{
            UserRepository: &database.UserRepository{
                SqlHandler: sqlHandler,
            },
        },
    }
}
*/

// GetUserByID は，httpを受け取り，portを組み立てて，inputPort.GetUserByIDを呼び出します．
func (u *User) GetUserByID(ctx *gin.Context) {

	//id, _ := strconv.Atoi(c.Param("id"))
	id := ctx.Param("id")

	outputPort := u.OutputFactory(ctx)
	repository := u.RepoFactory(u.DbHandler)
	inputPort := u.InputFactory(outputPort, repository)
	inputPort.GetUserByID(id)

}