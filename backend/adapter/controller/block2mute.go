package controller

import (
	"errors"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block2Mute struct {
	// -> presenter.NewBlock2MuteOutputPort
	OutputFactory func(
		contextHandler handler.ContextHandler,
		LoggerHandler handler.LoggerHandler,
	) port.Block2MuteOutputPort

	// -> interactor.NewBlock2MuteInputPort
	InputFactory func(
		o port.Block2MuteOutputPort,
		u port.Block2MuteRepository,
		LoggerHandler handler.LoggerHandler,
	) port.Block2MuteInputPort

	// -> gateway.NewBlock2MuteRepository
	RepoFactory func(
		LoggerHandler handler.LoggerHandler,
		dbHandler handler.DBConnectionHandler,
		twitterHandler handler.TwitterHandler,
		sessionHandler handler.SessionHandler,
	) port.Block2MuteRepository

	LoggerHandler  handler.LoggerHandler
	TwitterHandler handler.TwitterHandler
	DBHandler      handler.DBConnectionHandler
	SessionHandler handler.SessionHandler
}

func (b *Block2Mute) All(contextHandler handler.ContextHandler) {
	b.SessionHandler.SetContextHandler(contextHandler)

	id := b.SessionHandler.Get("user_id")

	outputPort := b.OutputFactory(contextHandler, b.LoggerHandler)
	repository := b.RepoFactory(b.LoggerHandler, b.DBHandler, b.TwitterHandler, b.SessionHandler)
	inputPort := b.InputFactory(outputPort, repository, b.LoggerHandler)

	if id == nil {
		outputPort.RenderError(errors.New("session user_id is not found"))
		return
	}
	inputPort.All(id.(string))
}
