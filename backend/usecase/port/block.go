package port

import (
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

type BlockInputPort interface {
	GetUserIDs(userID string)
}

type BlockOutputPort interface {
	Render(*entity.Blocks, int)
	RenderNotFound()
	RenderError(error)
}

type BlockRepository interface {
	GetUser(userID string) *entity.User
	GetBlocks(user *entity.User) (*entity.Blocks, int, error)
	TxUpdateAndDeleteBlocks(user *entity.User, blocks *entity.Blocks) error
}
