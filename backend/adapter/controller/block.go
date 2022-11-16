package controller

import (
	"errors"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block struct {
	// -> presenter.NewBlockOutputPort
	OutputFactory func(
		contextHandler handler.ContextHandler,
		LoggerHandler handler.LoggerHandler,
	) port.BlockOutputPort

	// -> interactor.NewBlockInputPort
	InputFactory func(
		o port.BlockOutputPort,
		u port.BlockRepository,
		LoggerHandler handler.LoggerHandler,
	) port.BlockInputPort

	// -> gateway.NewBlockRepository
	RepoFactory func(
		loggerHandler handler.LoggerHandler,
		dbHandler handler.DBHandler,
		twitterHandler handler.TwitterHandler,
		sessionHandler handler.SessionHandler,
	) port.BlockRepository

	LoggerHandler  handler.LoggerHandler
	TwitterHandler handler.TwitterHandler
	DBHandler      handler.DBHandler
	SessionHandler handler.SessionHandler
}

// GetBlockByID は，httpを受け取り，portを組み立てて，inputPort.GetBlockByIDを呼び出します．
func (b *Block) GetBlockByID(contextHandler handler.ContextHandler) {
	b.SessionHandler.SetContextHandler(contextHandler)

	id := b.SessionHandler.Get("user_id")

	outputPort := b.OutputFactory(contextHandler, b.LoggerHandler)
	repository := b.RepoFactory(b.LoggerHandler, b.DBHandler, b.TwitterHandler, b.SessionHandler)
	inputPort := b.InputFactory(outputPort, repository, b.LoggerHandler)

	if id == nil {
		outputPort.RenderError(errors.New("session user_id is not found"))
		return
	}
	inputPort.GetUserIDs(id.(string))
}
