package interactor

import (
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block struct {
	OutputPort port.BlockOutputPort
	BlockRepo  port.BlockRepository
}

func NewBlockInputPort(outputPort port.BlockOutputPort, blockRepository port.BlockRepository) port.BlockInputPort {
	return &Block{
		OutputPort: outputPort,
		BlockRepo:  blockRepository,
	}
}

func (b *Block) GetUserIDs(userID string) {
	blocks, total, err := b.BlockRepo.GetUserIDs(userID)
	if err != nil {
		b.OutputPort.RenderError(err)
		return
	}
	b.OutputPort.Render(blocks, total)
}
