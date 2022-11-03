package interactor

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type User struct {
	OutputPort    port.UserOutputPort
	UserRepo      port.UserRepository
	LoggerHandler handler.LoggerHandler
}

func NewUserInputPort(
	outputPort port.UserOutputPort,
	userRepository port.UserRepository,
	loggerHandler handler.LoggerHandler,
) port.UserInputPort {
	return &User{
		OutputPort:    outputPort,
		UserRepo:      userRepository,
		LoggerHandler: loggerHandler,
	}
}

func (u *User) GetUserByID(userID string) {
	user, err := u.UserRepo.GetUserByID(userID)
	if err != nil {
		u.OutputPort.RenderError(err)
		return
	}
	u.OutputPort.Render(user)
}
