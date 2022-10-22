package interactor

import (
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block2Mute struct {
	OutputPort     port.Block2MuteOutputPort
	Block2MuteRepo port.Block2MuteRepository
}

func NewBlock2MuteInputPort(outputPort port.Block2MuteOutputPort, block2MuteRepository port.Block2MuteRepository) port.Block2MuteInputPort {
	return &Block2Mute{
		OutputPort:     outputPort,
		Block2MuteRepo: block2MuteRepository,
	}
}

func (b *Block2Mute) All(userID string) {
	block2mute, err := b.Block2MuteRepo.All(userID)
	if err != nil {
		b.OutputPort.RenderError(err)
		return
	}
	b.OutputPort.Render(block2mute)
}
