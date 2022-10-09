package port

import (
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

type BlockInputPort interface {
	GetUserIDs(userID string)
}

type BlockOutputPort interface {
	Render(*[]entity.Block, int)
	RenderNotFound()
	RenderError(error)
}

type BlockRepository interface {
	GetUserIDs(userID string) (*[]entity.Block, int, error)
}
