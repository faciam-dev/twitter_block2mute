package port

import (
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

type UserInputPort interface {
	GetUserByID(userID string)
}

type UserOutputPort interface {
	Render(*entity.User)
	RenderNotFound()
	RenderError(error)
}

type UserRepository interface {
	GetUserByID(userID string) (*entity.User, error)
}