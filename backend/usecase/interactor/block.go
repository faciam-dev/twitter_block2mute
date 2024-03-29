package interactor

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block struct {
	OutputPort    port.BlockOutputPort
	BlockRepo     port.BlockRepository
	LoggerHandler handler.LoggerHandler
}

func NewBlockInputPort(
	outputPort port.BlockOutputPort,
	blockRepository port.BlockRepository,
	loggerHandler handler.LoggerHandler,
) port.BlockInputPort {
	return &Block{
		OutputPort:    outputPort,
		BlockRepo:     blockRepository,
		LoggerHandler: loggerHandler,
	}
}

func (b *Block) GetUserIDs(userID string) {
	user := b.BlockRepo.GetUser(userID)
	blocks, total, err := b.BlockRepo.GetBlocks(user)
	if err != nil {
		b.OutputPort.RenderError(err)
		return
	}
	err = b.BlockRepo.TxUpdateAndDeleteBlocks(user, blocks)
	if err != nil {
		b.OutputPort.RenderError(err)
		return
	}
	b.OutputPort.Render(blocks, total)
}
