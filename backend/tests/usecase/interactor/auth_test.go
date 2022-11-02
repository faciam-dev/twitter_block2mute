package interactor_test

import (
	"errors"
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	mock_port "github.com/faciam_dev/twitter_block2mute/backend/tests/usecase/port/mock"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/interactor"
	"github.com/golang/mock/gomock"
)

// for Auth / IsAuthenticated
type argsAuth struct {
	Authenticated   int
	AuthUrl         string
	RepositoryError error
}

var (
	tableAuth = []struct {
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
)

// for callback
type argsCallback struct {
	Authenticated   int
	AuthUrl         string
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
				Authenticated:   0,
				AuthUrl:         "http://testurl",
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
				AuthUrl:         "http://testurl",
				RepositoryError: errors.New("user is not found"),
				Token:           "errortoken",
				Secret:          "errorsecret",
				TwitterID:       "1234567890",
				TwitterName:     "test",
			},
		},
	}
)

// for Logout
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

// Auth()のテスト
func TestAuth(t *testing.T) {
	for _, tt := range tableAuth {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// エンティティ類
			fromRepositoryAuth := &entity.Auth{Authenticated: tt.args.Authenticated, AuthUrl: tt.args.AuthUrl}

			// outputPort の設定
			outputPort := mock_port.NewMockAuthOutputPort(mockCtrl)
			outputPort.EXPECT().RenderAuth(fromRepositoryAuth).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// authRepositoryモックの設定
			authRepository := mock_port.NewMockAuthRepository(mockCtrl)
			authRepository.EXPECT().Auth().Return(fromRepositoryAuth, tt.args.RepositoryError)

			// テスト対象の構築
			interactor := interactor.NewAuthInputPort(outputPort, authRepository)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.Auth()
		})
	}
}

// IsAuthenticated()のテスト
func TestIsAuthenticated(t *testing.T) {
	for _, tt := range tableAuth {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// エンティティ類
			fromRepositoryAuth := &entity.Auth{Authenticated: tt.args.Authenticated, AuthUrl: tt.args.AuthUrl}

			// outputPort の設定
			outputPort := mock_port.NewMockAuthOutputPort(mockCtrl)
			outputPort.EXPECT().RenderIsAuth(fromRepositoryAuth).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// authRepositoryモックの設定
			authRepository := mock_port.NewMockAuthRepository(mockCtrl)
			authRepository.EXPECT().IsAuthenticated().Return(fromRepositoryAuth, tt.args.RepositoryError)

			// テスト対象の構築
			interactor := interactor.NewAuthInputPort(outputPort, authRepository)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.IsAuthenticated()
		})
	}
}

// Callback()のテスト
func TestCallback(t *testing.T) {
	for _, tt := range tableCallback {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// エンティティ類
			fromRepositoryAuth := &entity.Auth{Authenticated: tt.args.Authenticated, AuthUrl: tt.args.AuthUrl}

			// outputPort の設定
			outputPort := mock_port.NewMockAuthOutputPort(mockCtrl)
			outputPort.EXPECT().RenderCallback(fromRepositoryAuth).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// authRepositoryモックの設定
			authRepository := mock_port.NewMockAuthRepository(mockCtrl)
			authRepository.EXPECT().Callback(
				tt.args.Token,
				tt.args.Secret,
				tt.args.TwitterID,
				tt.args.TwitterName,
			).Return(fromRepositoryAuth, tt.args.RepositoryError)

			// テスト対象の構築
			interactor := interactor.NewAuthInputPort(outputPort, authRepository)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.Callback(
				tt.args.Token,
				tt.args.Secret,
				tt.args.TwitterID,
				tt.args.TwitterName,
			)
		})
	}
}

func TestLogout(t *testing.T) {
	for _, tt := range tableLogout {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// エンティティ類
			fromRepositoryAuth := &entity.Auth{Logout: tt.args.Logout}

			// outputPort の設定
			outputPort := mock_port.NewMockAuthOutputPort(mockCtrl)
			outputPort.EXPECT().RenderLogout(fromRepositoryAuth).Return().AnyTimes()
			outputPort.EXPECT().RenderError(tt.args.RepositoryError).Return().AnyTimes()

			// authRepositoryモックの設定
			authRepository := mock_port.NewMockAuthRepository(mockCtrl)
			authRepository.EXPECT().Logout().Return(fromRepositoryAuth, tt.args.RepositoryError)

			// テスト対象の構築
			interactor := interactor.NewAuthInputPort(outputPort, authRepository)

			// 得られるものはないため、実行できればテスト通過とする。
			interactor.Logout()
		})
	}
}
