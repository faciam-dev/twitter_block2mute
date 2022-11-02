package controller

import (
	"errors"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block2Mute struct {
	// -> presenter.NewBlock2MuteOutputPort
	OutputFactory func(contextHandler handler.ContextHandler) port.Block2MuteOutputPort
	// -> interactor.NewBlock2MuteInputPort
	InputFactory func(o port.Block2MuteOutputPort, u port.Block2MuteRepository) port.Block2MuteInputPort
	// -> gateway.NewBlock2MuteRepository
	RepoFactory func(
		dbHandler handler.BlockDbHandler,
		userDbHandler handler.UserDbHandler,
		muteDbHandler handler.MuteDbHandler,
		twitterHandler handler.TwitterHandler,
		sessionHandler handler.SessionHandler) port.Block2MuteRepository
	TwitterHandler handler.TwitterHandler
	BlockDbHandler handler.BlockDbHandler
	MuteDbHandler  handler.MuteDbHandler
	UserDbHandler  handler.UserDbHandler
	SessionHandler handler.SessionHandler
}

// GetBlockByID は，httpを受け取り，portを組み立てて，inputPort.GetBlockByIDを呼び出します．
func (b *Block2Mute) All(contextHandler handler.ContextHandler) {
	b.SessionHandler.SetContextHandler(contextHandler)

	id := b.SessionHandler.Get("user_id")

	outputPort := b.OutputFactory(contextHandler)
	repository := b.RepoFactory(b.BlockDbHandler, b.UserDbHandler, b.MuteDbHandler, b.TwitterHandler, b.SessionHandler)
	inputPort := b.InputFactory(outputPort, repository)

	if id == nil {
		outputPort.RenderError(errors.New("session user_id is not found"))
		return
	}
	inputPort.All(id.(string))
}
