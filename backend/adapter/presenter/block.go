package presenter

import (
	"net/http"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block struct {
	contextHandler handler.ContextHandler
	loggerHandler  handler.LoggerHandler
}

// NewBlockOutputPort はBlockOutputPortを取得します．
func NewBlockOutputPort(
	contextHandler handler.ContextHandler,
	loggerHandler handler.LoggerHandler,
) port.BlockOutputPort {
	return &Block{
		contextHandler: contextHandler,
		loggerHandler:  loggerHandler,
	}
}

// usecase.BlockOutputPortを実装している
// Render はBlockモデルを出力します．
func (u *Block) Render(blocks *[]entity.Block, total int) {
	twitterIDs := []string{}
	for _, block := range *blocks {
		twitterIDs = append(twitterIDs, block.TargetTwitterID)
	}

	u.contextHandler.JSON(http.StatusOK, map[string]interface{}{
		"total": total,
		"ids":   twitterIDs,
	})
}

// RenderNotFound は ないことを出力します。
func (u *Block) RenderNotFound() {
	u.contextHandler.JSON(http.StatusNotFound, map[string]interface{}{})
}

// RenderError はErrorを出力します．
func (u *Block) RenderError(err error) {
	u.contextHandler.JSON(http.StatusInternalServerError, map[string]interface{}{
		"error": err,
	})
}
