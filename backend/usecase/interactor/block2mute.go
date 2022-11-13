package interactor

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block2Mute struct {
	OutputPort     port.Block2MuteOutputPort
	Block2MuteRepo port.Block2MuteRepository
	LoggerHandler  handler.LoggerHandler
}

func NewBlock2MuteInputPort(
	outputPort port.Block2MuteOutputPort,
	block2MuteRepository port.Block2MuteRepository,
	loggerHandler handler.LoggerHandler,
) port.Block2MuteInputPort {
	return &Block2Mute{
		OutputPort:     outputPort,
		Block2MuteRepo: block2MuteRepository,
		LoggerHandler:  loggerHandler,
	}
}

func (b *Block2Mute) All(userID string) {
	user := b.Block2MuteRepo.GetUser(userID)
	if err := b.Block2MuteRepo.AuthTwitter(); err != nil {
		b.OutputPort.RenderError(err)
		return
	}
	block2mute, err := b.Block2MuteRepo.All(user)
	if err != nil {
		b.OutputPort.RenderError(err)
		return
	}
	b.OutputPort.Render(block2mute)
}
