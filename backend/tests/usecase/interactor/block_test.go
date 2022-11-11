package interactor_test

import (
	"errors"
	"testing"
	"time"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/tests/adapter/gateway/mock_handler"
	mock_port "github.com/faciam_dev/twitter_block2mute/backend/tests/usecase/port/mock"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/interactor"
	"github.com/golang/mock/gomock"
)

// GetUserIDsのテスト
func TestGetUserIDs(t *testing.T) {
	type argsGetUserIDs struct {
		UserID          string
		Total           int
		BlockEntity     entity.Blocks
		RepositoryError error
	}
	nowTime := time.Now()
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
					BlockEntity: entity.Blocks{
						*entity.NewBlock(1, "2", 0, nowTime, nowTime),
						*entity.NewBlock(2, "3", 0, nowTime, nowTime),
						*entity.NewBlock(3, "4", 0, nowTime, nowTime),
					},
					RepositoryError: nil,
				},
			},
			{
				name: "error",
				args: argsGetUserIDs{
					UserID:          "20000",
					Total:           0,
					BlockEntity:     entity.Blocks{},
					RepositoryError: errors.New("blocks are not found"),
				},
			},
		}
	)

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
			outputPort.EXPECT().Render(&tt.args.BlockEntity, tt.args.Total).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// userRepositoryモックの設定
			blockRepository := mock_port.NewMockBlockRepository(mockCtrl)
			blockRepository.EXPECT().GetUser(tt.args.UserID).Return(entity.NewBlankUser())
			blockRepository.EXPECT().GetBlocks(entity.NewBlankUser()).Return(&tt.args.BlockEntity, tt.args.Total, tt.args.RepositoryError)
			blockRepository.EXPECT().TxUpdateAndDeleteBlocks(entity.NewBlankUser(), &tt.args.BlockEntity).Return(tt.args.RepositoryError).AnyTimes()

			// テスト対象の構築
			interactor := interactor.NewBlockInputPort(outputPort, blockRepository, logger)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.GetUserIDs(tt.args.UserID)
		})
	}
}
