package controller

import (
	"errors"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block struct {
	// -> presenter.NewBlockOutputPort
	OutputFactory func(contextHandler handler.ContextHandler) port.BlockOutputPort
	// -> interactor.NewBlockInputPort
	InputFactory func(o port.BlockOutputPort, u port.BlockRepository) port.BlockInputPort
	// -> gateway.NewBlockRepository
	RepoFactory    func(dbHandler handler.BlockDbHandler, userDbHandler handler.UserDbHandler, twitterHandler handler.TwitterHandler, sessionHandler handler.SessionHandler) port.BlockRepository
	TwitterHandler handler.TwitterHandler
	BlockDbHandler handler.BlockDbHandler
	UserDbHandler  handler.UserDbHandler
	SessionHandler handler.SessionHandler
}

// GetBlockByID は，httpを受け取り，portを組み立てて，inputPort.GetBlockByIDを呼び出します．
func (b *Block) GetBlockByID(contextHandler handler.ContextHandler) {
	b.SessionHandler.SetContextHandler(contextHandler)

	id := b.SessionHandler.Get("user_id")

	outputPort := b.OutputFactory(contextHandler)
	repository := b.RepoFactory(b.BlockDbHandler, b.UserDbHandler, b.TwitterHandler, b.SessionHandler)
	inputPort := b.InputFactory(outputPort, repository)

	if id == nil {
		outputPort.RenderError(errors.New("session user_id is not found"))
		return
	}
	inputPort.GetUserIDs(id.(string))
}
