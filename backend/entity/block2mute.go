package entity

type Block2Mute struct {
	numberOfSuccess   uint
	successTwitterIDs []string
}

func NewBlock2Mute(
	numberOfSuccess uint,
	successTwitterIDs []string,
) *Block2Mute {
	block2mute := Block2Mute{}
	block2mute.numberOfSuccess = numberOfSuccess
	block2mute.successTwitterIDs = successTwitterIDs

	return &block2mute
}

func (b *Block2Mute) GetNumberOfSuccess() uint {
	return b.numberOfSuccess
}

func (b *Block2Mute) GetSuccessTwitterIDs() []string {
	return b.successTwitterIDs
}
