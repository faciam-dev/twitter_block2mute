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

// Auth()のテスト
func TestAuth(t *testing.T) {
	type argsAuth struct {
		Authenticated   int
		AuthUrl         string
		RepositoryError error
	}
	tableAuth := []struct {
		name string
		args argsAuth
		err  error
	}{
		{
			name: "success",
			args: argsAuth{
				Authenticated:   0,
				AuthUrl:         "http://testurl",
				RepositoryError: nil,
			},
		},
		{
			name: "error",
			args: argsAuth{
				Authenticated:   0,
				AuthUrl:         "http://testurl",
				RepositoryError: errors.New("user is not found"),
			},
		},
	}
	for _, tt := range tableAuth {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// エンティティ類
			auth := entity.NewAuth(tt.args.AuthUrl)
			if tt.args.Authenticated == 1 {
				auth.SuccessAuthenticated()
			}
			fromRepositoryAuth := auth

			// loggerモックの設定
			logger := mock_handler.NewMockLoggerHandler(mockCtrl)
			logger.EXPECT().Debugf(gomock.Any()).AnyTimes().Return()

			// outputPort の設定
			outputPort := mock_port.NewMockAuthOutputPort(mockCtrl)
			outputPort.EXPECT().RenderAuth(fromRepositoryAuth).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// authRepositoryモックの設定
			authRepository := mock_port.NewMockAuthRepository(mockCtrl)
			authRepository.EXPECT().GetAuthUrl().Return(tt.args.AuthUrl, tt.args.RepositoryError)

			// テスト対象の構築
			interactor := interactor.NewAuthInputPort(outputPort, authRepository, logger)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.Auth()
		})
	}
}

// IsAuthenticated()のテスト
func TestIsAuthenticated(t *testing.T) {
	type argsAuthenticated struct {
		Authenticated   int
		RepositoryError error
	}
	tableAuthenticated := []struct {
		name string
		args argsAuthenticated
		err  error
	}{
		{
			name: "success",
			args: argsAuthenticated{
				Authenticated:   1,
				RepositoryError: nil,
			},
		},
		{
			name: "error",
			args: argsAuthenticated{
				Authenticated:   0,
				RepositoryError: errors.New("user is not found"),
			},
		},
	}
	for _, tt := range tableAuthenticated {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// エンティティ類
			auth := entity.NewAuth("")
			if tt.args.Authenticated == 1 {
				auth.SuccessAuthenticated()
			}
			fromRepositoryAuth := auth

			// loggerモックの設定
			logger := mock_handler.NewMockLoggerHandler(mockCtrl)
			logger.EXPECT().Debugf(gomock.Any()).AnyTimes().Return()
			logger.EXPECT().Errorw(gomock.Any(), gomock.Any()).AnyTimes().Return()

			// outputPort の設定
			outputPort := mock_port.NewMockAuthOutputPort(mockCtrl)
			outputPort.EXPECT().RenderIsAuth(fromRepositoryAuth).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// authRepositoryモックの設定
			authRepository := mock_port.NewMockAuthRepository(mockCtrl)
			authRepository.EXPECT().IsAuthenticated().Return(tt.args.RepositoryError)

			// テスト対象の構築
			interactor := interactor.NewAuthInputPort(outputPort, authRepository, logger)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.IsAuthenticated()
		})
	}
}

// Callback()のテスト
func TestCallback(t *testing.T) {
	// for callback
	type argsCallback struct {
		Authenticated   int
		RepositoryError error
		Token           string
		Secret          string
		TwitterID       string
		TwitterName     string
	}

	var (
		tableCallback = []struct {
			name string
			args argsCallback
			err  error
		}{
			{
				name: "success",
				args: argsCallback{
					Authenticated:   1,
					RepositoryError: nil,
					Token:           "token",
					Secret:          "secret",
					TwitterID:       "1234567890",
					TwitterName:     "test",
				},
			},
			{
				name: "error",
				args: argsCallback{
					Authenticated:   0,
					RepositoryError: errors.New("user is not found"),
					Token:           "errortoken",
					Secret:          "errorsecret",
					TwitterID:       "1234567890",
					TwitterName:     "test",
				},
			},
		}
	)
	for _, tt := range tableCallback {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// エンティティ類
			auth := entity.NewAuth("")
			if tt.args.Authenticated == 1 {
				auth.SuccessAuthenticated()
			}
			fromRepositoryAuth := auth

			// loggerモックの設定
			logger := mock_handler.NewMockLoggerHandler(mockCtrl)
			logger.EXPECT().Debugf(gomock.Any()).AnyTimes().Return()

			// outputPort の設定
			outputPort := mock_port.NewMockAuthOutputPort(mockCtrl)
			outputPort.EXPECT().RenderCallback(fromRepositoryAuth).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// authRepositoryの準備
			mockTwitterCredentials := mock_handler.NewMockTwitterCredentials(mockCtrl)
			mockTwitterCredentials.EXPECT().GetToken().Return(tt.args.Token).AnyTimes()
			mockTwitterCredentials.EXPECT().GetSecret().Return(tt.args.Secret).AnyTimes()

			mockTwitterValue := mock_handler.NewMockTwitterValues(mockCtrl)
			mockTwitterValue.EXPECT().GetTwitterID().Return(tt.args.TwitterID).AnyTimes()
			mockTwitterValue.EXPECT().GetTwitterScreenName().Return(tt.args.TwitterName).AnyTimes()

			user := &entity.User{} // TODO: 次のPRでファクトリ関数による生成にする。

			// authRepositoryモックの設定
			authRepository := mock_port.NewMockAuthRepository(mockCtrl)
			authRepository.EXPECT().AuthByCallbackParams(
				tt.args.Token,
				tt.args.Secret,
			).Return(
				mockTwitterCredentials,
				mockTwitterValue,
				tt.args.RepositoryError,
			).AnyTimes()

			authRepository.EXPECT().FindUserByTwitterID(
				mockTwitterValue.GetTwitterID(),
			).Return(
				user,
				tt.args.RepositoryError,
			).AnyTimes()

			authRepository.EXPECT().UpsertUser(user).Return(tt.args.RepositoryError).AnyTimes()

			authRepository.EXPECT().UpdateTwitterApi(
				mockTwitterCredentials.GetToken(),
				mockTwitterCredentials.GetSecret(),
			).Return().AnyTimes()

			authRepository.EXPECT().UpdateSession(
				mockTwitterCredentials.GetToken(),
				mockTwitterCredentials.GetSecret(),
				int(user.ID),
				mockTwitterValue.GetTwitterID(),
			).Return(tt.args.RepositoryError).AnyTimes()

			// テスト対象の構築
			interactor := interactor.NewAuthInputPort(outputPort, authRepository, logger)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.Callback(
				tt.args.Token,
				tt.args.Secret,
			)
		})
	}
}

func TestLogout(t *testing.T) {
	type argsLogout struct {
		Logout          int
		RepositoryError error
	}

	var (
		tableLogout = []struct {
			name string
			args argsLogout
			err  error
		}{
			{
				name: "success",
				args: argsLogout{
					Logout:          1,
					RepositoryError: nil,
				},
			},
			{
				name: "error",
				args: argsLogout{
					Logout:          0,
					RepositoryError: errors.New(" is not found"),
				},
			},
		}
	)
	for _, tt := range tableLogout {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// エンティティ類
			auth := entity.NewAuth("")
			if tt.args.Logout == 1 {
				auth.SuccessLogout()
			}
			fromRepositoryAuth := auth

			// loggerモックの設定
			logger := mock_handler.NewMockLoggerHandler(mockCtrl)
			logger.EXPECT().Debugf(gomock.Any()).AnyTimes().Return()

			// outputPort の設定
			outputPort := mock_port.NewMockAuthOutputPort(mockCtrl)
			outputPort.EXPECT().RenderLogout(fromRepositoryAuth).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// authRepositoryモックの設定
			authRepository := mock_port.NewMockAuthRepository(mockCtrl)
			authRepository.EXPECT().Logout().Return(tt.args.RepositoryError)

			// テスト対象の構築
			interactor := interactor.NewAuthInputPort(outputPort, authRepository, logger)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.Logout()
		})
	}
}
