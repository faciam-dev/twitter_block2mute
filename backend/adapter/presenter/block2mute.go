package presenter

import (
	"net/http"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block2Mute struct {
	contextHandler handler.ContextHandler
	loggerHandler  handler.LoggerHandler
}

// NewBlockOutputPort はBlockOutputPortを取得します．
func NewBlock2MuteOutputPort(
	contextHandler handler.ContextHandler,
	loggerHandler handler.LoggerHandler,
) port.Block2MuteOutputPort {
	return &Block2Mute{
		contextHandler: contextHandler,
		loggerHandler:  loggerHandler,
	}
}

// usecase.BlockOutputPortを実装している
// Render はBlockモデルを出力します．
func (b *Block2Mute) Render(block2Mute *entity.Block2Mute) {
	b.contextHandler.JSON(http.StatusOK, map[string]interface{}{
		"num_success": block2Mute.NumberOfSuccess,
		"success_ids": block2Mute.SuccessTwitterIDs,
	})
}

// RenderNotFound は ないことを出力します。
func (b *Block2Mute) RenderNotFound() {
	b.contextHandler.JSON(http.StatusNotFound, map[string]interface{}{})
}

// RenderError はErrorを出力します．
func (b *Block2Mute) RenderError(err error) {
	b.contextHandler.JSON(http.StatusInternalServerError, map[string]interface{}{
		"error": err,
	})
}
