package interactor_test

import (
	"errors"
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	mock_port "github.com/faciam_dev/twitter_block2mute/backend/tests/usecase/port/mock"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/interactor"
	"github.com/golang/mock/gomock"
)

// for GetUserIDs
type argsAll struct {
	UserID          string
	Block2Mute      entity.Block2Mute
	RepositoryError error
}

var (
	tableAll = []struct {
		name string
		args argsAll
		err  error
	}{
		{
			name: "success",
			args: argsAll{
				UserID: "1",
				Block2Mute: entity.Block2Mute{
					NumberOfSuccess:   1,
					SuccessTwitterIDs: []string{"1"},
				},
				RepositoryError: nil,
			},
		},
		{
			name: "error",
			args: argsAll{
				UserID:          "20000",
				Block2Mute:      entity.Block2Mute{},
				RepositoryError: errors.New("blocks are not found"),
			},
		},
	}
)

// Allのテスト
func TestAll(t *testing.T) {
	for _, tt := range tableAll {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// outputPort の設定
			outputPort := mock_port.NewMockBlock2MuteOutputPort(mockCtrl)
			outputPort.EXPECT().Render(&tt.args.Block2Mute).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// userRepositoryモックの設定
			block2MuteRepository := mock_port.NewMockBlock2MuteRepository(mockCtrl)
			block2MuteRepository.EXPECT().All(tt.args.UserID).Return(&tt.args.Block2Mute, tt.args.RepositoryError)

			// テスト対象の構築
			interactor := interactor.NewBlock2MuteInputPort(outputPort, block2MuteRepository)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.All(tt.args.UserID)
		})
	}
}
