package interactor_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	mock_port "github.com/faciam_dev/twitter_block2mute/backend/tests/usecase/port/mock"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/interactor"
	"github.com/golang/mock/gomock"
)

// for GetUserByID
type argsGetUserByID struct {
	UserID          string
	Authenticated   int
	AuthUrl         string
	RepositoryError error
}

var (
	tableGetUserByID = []struct {
		name string
		args argsGetUserByID
		err  error
	}{
		{
			name: "success",
			args: argsGetUserByID{
				UserID:          "1",
				RepositoryError: nil,
			},
		},
		{
			name: "error",
			args: argsGetUserByID{
				UserID:          "20000",
				RepositoryError: errors.New("user is not found"),
			},
		},
	}
)

// GetUserByID()のテスト
func TestGetUserByID(t *testing.T) {
	for _, tt := range tableGetUserByID {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// エンティティ類
			userID, _ := strconv.ParseUint(tt.args.UserID, 10, 64)
			fromRepositoryUser := &entity.User{ID: uint(userID)}

			// outputPort の設定
			outputPort := mock_port.NewMockUserOutputPort(mockCtrl)
			outputPort.EXPECT().Render(fromRepositoryUser).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// userRepositoryモックの設定
			userRepository := mock_port.NewMockUserRepository(mockCtrl)
			userRepository.EXPECT().GetUserByID(tt.args.UserID).Return(fromRepositoryUser, tt.args.RepositoryError)

			// テスト対象の構築
			interactor := interactor.NewUserInputPort(outputPort, userRepository)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.GetUserByID(tt.args.UserID)
		})
	}
}
