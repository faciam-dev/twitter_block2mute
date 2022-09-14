package controller

import (
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	// -> presenter.NewUserOutputPort
	OutputFactory func(c *gin.Context) port.UserOutputPort
	// -> interactor.NewUserInputPort
	InputFactory func(o port.UserOutputPort, u port.UserRepository) port.UserInputPort
	// -> gateway.NewUserRepository
	RepoFactory func(c *gorm.DB) port.UserRepository
	Conn *gorm.DB
}

// GetUserByID は，httpを受け取り，portを組み立てて，inputPort.GetUserByIDを呼び出します．
func (u *User) GetUserByID(c *gin.Context) {

	//id, _ := strconv.Atoi(c.Param("id"))
	id := c.Param("id")

	outputPort := u.OutputFactory(c)
	repository := u.RepoFactory(u.Conn)
	inputPort := u.InputFactory(outputPort, repository)
	inputPort.GetUserByID(id)

}