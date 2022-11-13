package port

import (
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

type Block2MuteInputPort interface {
	All(userID string)
}

type Block2MuteOutputPort interface {
	Render(*entity.Block2Mute)
	RenderNotFound()
	RenderError(error)
}

type Block2MuteRepository interface {
	GetUser(userID string) *entity.User
	AuthTwitter() error
	All(user *entity.User) (*entity.Block2Mute, error)
}
