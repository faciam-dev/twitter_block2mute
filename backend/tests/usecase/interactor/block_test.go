package interactor_test

import (
	"errors"
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/tests/adapter/gateway/mock_handler"
	mock_port "github.com/faciam_dev/twitter_block2mute/backend/tests/usecase/port/mock"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/interactor"
	"github.com/golang/mock/gomock"
)

// for GetUserIDs
type argsGetUserIDs struct {
	UserID          string
	Total           int
	BlockEntities   []entity.Block
	RepositoryError error
}

var (
	tableGetUserIDs = []struct {
		name string
		args argsGetUserIDs
		err  error
	}{
		{
			name: "success",
			args: argsGetUserIDs{
				UserID: "1",
				Total:  3,
				BlockEntities: []entity.Block{
					{TargetTwitterID: "2"},
					{TargetTwitterID: "3"},
					{TargetTwitterID: "4"},
				},
				RepositoryError: nil,
			},
		},
		{
			name: "error",
			args: argsGetUserIDs{
				UserID:          "20000",
				Total:           0,
				BlockEntities:   []entity.Block{},
				RepositoryError: errors.New("blocks are not found"),
			},
		},
	}
)

// GetUserIDsのテスト
func TestGetUserIDs(t *testing.T) {
	for _, tt := range tableGetUserIDs {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// loggerモックの設定
			logger := mock_handler.NewMockLoggerHandler(mockCtrl)
			logger.EXPECT().Debugf(gomock.Any()).AnyTimes().Return()

			// outputPort の設定
			outputPort := mock_port.NewMockBlockOutputPort(mockCtrl)
			outputPort.EXPECT().Render(&tt.args.BlockEntities, tt.args.Total).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// userRepositoryモックの設定
			blockRepository := mock_port.NewMockBlockRepository(mockCtrl)
			blockRepository.EXPECT().GetUserIDs(tt.args.UserID).Return(&tt.args.BlockEntities, tt.args.Total, tt.args.RepositoryError)

			// テスト対象の構築
			interactor := interactor.NewBlockInputPort(outputPort, blockRepository, logger)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.GetUserIDs(tt.args.UserID)
		})
	}
}
