package gateway_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/common"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/tests/adapter/gateway/mock_handler"
)

// IsAuthnticatedに対するテスト
func TestIsAuthenticated(t *testing.T) {
	type args struct {
		SessionToken       interface{}
		SessionTokenSecret interface{}
		TwitterName        string
		TwitterAccountName string
		TwitterID          string
		OAuthToken         string
		OAuthTokenSecret   string
	}
	table := []struct {
		name string
		args args
		want entity.Auth
		err  error
	}{
		{
			name: "success_1",
			args: args{
				SessionToken:       "token",
				SessionTokenSecret: "secret",
				TwitterName:        "test",
				TwitterAccountName: "Test",
				TwitterID:          "1234567890",
				OAuthToken:         "token",
				OAuthTokenSecret:   "secret",
			},
			want: *entity.NewAuth("").SuccessAuthenticated(),
			err:  nil,
		},
		{
			name: "success_2",
			args: args{
				SessionToken:       nil,
				SessionTokenSecret: nil,
				TwitterName:        "test",
				TwitterAccountName: "Test",
				TwitterID:          "1234567890",
				OAuthToken:         "token",
				OAuthTokenSecret:   "secret",
			},
			want: *entity.NewAuth(""),
			err:  nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			// モックの生成
			// sqlmock処理
			dbUserHandler /*dbMock*/, _, err := newMockGormDBUserHandler()

			if err != nil {
				t.Error("sqlmock not work")
			}

			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// contextHandler
			contextHandler := mock_handler.NewMockContextHandler(mockCtrl)

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes().Return()

			// twtterHandler
			mockTwitterUser := mock_handler.NewMockTwitterUser(mockCtrl)
			mockTwitterUser.EXPECT().GetTwitterID().Return(args.TwitterID).AnyTimes()
			mockTwitterUser.EXPECT().GetTwitterScreenName().Return(args.TwitterName).AnyTimes()
			mockTwitterUser.EXPECT().GetTwitterName().Return(args.TwitterAccountName).AnyTimes()

			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)
			mockTwitterHandler.EXPECT().GetUser(args.TwitterID).Return(mockTwitterUser, nil).AnyTimes()
			mockTwitterHandler.EXPECT().UpdateTwitterApi(args.OAuthToken, args.OAuthTokenSecret).Return().AnyTimes()

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().SetContextHandler(contextHandler).Return()
			sessionHandler.EXPECT().Get("token").Return(args.SessionToken).AnyTimes()
			sessionHandler.EXPECT().Get("secret").Return(args.SessionTokenSecret).AnyTimes()
			sessionHandler.EXPECT().Get("twitter_id").Return(args.TwitterID).AnyTimes()

			// repository
			authRepository := gateway.NewAuthRepository(
				contextHandler,
				loggerHandler,
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			err = authRepository.IsAuthenticated()
			if tt.err != err {
				t.Errorf("IsAuthenticated() err = %v, want %v", err, tt.err)
			}
		})
	}
}

// Authに対するテスト
func TestAuth(t *testing.T) {
	type args struct {
		AuthUrl string
	}
	table := []struct {
		name string
		args args
		want entity.Auth
		err  error
	}{
		{
			name: "success",
			args: args{
				AuthUrl: "https://api.twitter.com/oauth/authenticate",
			},
			want: *entity.NewAuth("https://api.twitter.com/oauth/authenticate"),
			err:  nil,
		},
	}

	dbUserHandler, _, err := newMockGormDBUserHandler()

	if err != nil {
		t.Error("sqlmock not work")
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// contextHandler
			contextHandler := mock_handler.NewMockContextHandler(mockCtrl)

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes().Return()

			// twtterHandler
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)
			mockTwitterHandler.EXPECT().AuthorizationURL().Return(args.AuthUrl+"?token="+common.RandomString(16), nil)

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().SetContextHandler(contextHandler).Return()

			authRepository := gateway.NewAuthRepository(
				contextHandler,
				loggerHandler,
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			got, err := authRepository.GetAuthUrl()

			if tt.err != err {
				t.Errorf("Auth() err = %v, want %v", err, tt.err)
			}

			if strings.Index(tt.want.GetAuthUrl(), got) > 0 {
				t.Errorf("Auth() auth = %v, want %v", got, tt.want.GetAuthUrl())
			}
		})
	}
}

// AuthByCallbackParamsに対するテスト
func TestAuthByCallbackParams(t *testing.T) {
	type args struct {
		SessionToken       interface{}
		SessionTokenSecret interface{}
		TwitterName        string
		TwitterAccountName string
		TwitterID          string
		OAuthToken         string
		OAuthTokenSecret   string
	}
	type wantResult struct {
		OAuthToken       string
		OAuthTokenSecret string
		TwitterName      string
		TwitterID        string
	}
	table := []struct {
		name string
		args args
		want wantResult
		err  error
	}{
		{
			name: "success",
			args: args{
				SessionToken:       "token",
				SessionTokenSecret: "secret",
				TwitterName:        "test",
				TwitterAccountName: "Test",
				TwitterID:          "1234567890",
				OAuthToken:         "token",
				OAuthTokenSecret:   "secret",
			},
			want: wantResult{
				OAuthToken:       "token",
				OAuthTokenSecret: "secret",
				TwitterName:      "test",
				TwitterID:        "1234567890",
			},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			dbUserHandler, _, err := newMockGormDBUserHandler()

			if err != nil {
				t.Error("sqlmock not work")
			}

			// モックの生成
			// gomock処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// contextHandler
			contextHandler := mock_handler.NewMockContextHandler(mockCtrl)

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)
			loggerHandler.EXPECT().Debug(gomock.Any()).AnyTimes().Return()
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes().Return()

			// twtterHandler
			mockTwitterCredentials := mock_handler.NewMockTwitterCredentials(mockCtrl)
			mockTwitterCredentials.EXPECT().GetToken().Return(args.SessionToken).AnyTimes()
			mockTwitterCredentials.EXPECT().GetSecret().Return(args.SessionTokenSecret).AnyTimes()

			mockTwitterValue := mock_handler.NewMockTwitterValues(mockCtrl)
			mockTwitterValue.EXPECT().GetTwitterID().Return(args.TwitterID).AnyTimes()
			mockTwitterValue.EXPECT().GetTwitterScreenName().Return(args.TwitterName).AnyTimes()

			mockTwitterUser := mock_handler.NewMockTwitterUser(mockCtrl)
			mockTwitterUser.EXPECT().GetTwitterID().Return(args.TwitterID).AnyTimes()
			mockTwitterUser.EXPECT().GetTwitterScreenName().Return(args.TwitterName).AnyTimes()
			mockTwitterUser.EXPECT().GetTwitterName().Return(args.TwitterAccountName).AnyTimes()

			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)
			mockTwitterHandler.EXPECT().GetUser(args.TwitterID).Return(mockTwitterUser, nil).AnyTimes()
			mockTwitterHandler.EXPECT().UpdateTwitterApi(args.OAuthToken, args.OAuthTokenSecret).Return().AnyTimes()
			mockTwitterHandler.EXPECT().GetCredentials(gomock.Any(), gomock.Any()).Return(mockTwitterCredentials, mockTwitterValue, nil)

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().SetContextHandler(contextHandler).Return()
			sessionHandler.EXPECT().Set(gomock.Any(), gomock.Any()).Return().AnyTimes()
			sessionHandler.EXPECT().Set(gomock.Any(), gomock.Any()).Return().AnyTimes()
			sessionHandler.EXPECT().Save().Return(nil).AnyTimes()

			// repository
			authRepository := gateway.NewAuthRepository(
				contextHandler,
				loggerHandler,
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			gotTwitterCredentials, gotTwitterValues, err := authRepository.AuthByCallbackParams(args.OAuthToken, args.OAuthTokenSecret)

			if tt.err != err {
				t.Errorf("AuthByCallbackParams() err = %v, want %v", err, tt.err)
			}

			if gotTwitterCredentials.GetSecret() != tt.want.OAuthTokenSecret || gotTwitterCredentials.GetToken() != tt.want.OAuthToken {
				t.Errorf(
					"AuthByCallbackParams() Credentials = secret:%v token:%v , want = secret: %v token:%v",
					gotTwitterCredentials.GetSecret(),
					gotTwitterCredentials.GetToken(),
					tt.want.OAuthTokenSecret,
					tt.want.OAuthToken,
				)
			}

			if gotTwitterValues.GetTwitterID() != tt.args.TwitterID || gotTwitterValues.GetTwitterScreenName() != tt.args.TwitterName {
				t.Errorf(
					"AuthByCallbackParams() TwitterValues = twitterID:%v twitterScreenName:%v , want = twitterID:%v twitterScreenName:%v",
					gotTwitterValues.GetTwitterID(),
					gotTwitterValues.GetTwitterScreenName(),
					tt.want.TwitterID,
					tt.want.TwitterName,
				)
			}

		})
	}
}

// FindUserByTwitterIDに対するテスト
func TestFindUserByTwitterID(t *testing.T) {
	type args struct {
		UserID             uint
		TwitterName        string
		TwitterAccountName string
		TwitterID          string
	}
	table := []struct {
		name string
		args args
		want *entity.User
		err  error
	}{
		{
			name: "success",
			args: args{
				UserID:             1,
				TwitterName:        "test",
				TwitterAccountName: "Test",
				TwitterID:          "1234567890",
			},
			want: entity.NewBlankUser().Update(1, "test", "Test", "1234567890"),
			err:  nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			dbUserHandler, dbMock, err := newMockGormDBUserHandler()

			if err != nil {
				t.Error("sqlmock not work")
			}

			// モックの生成
			// sqlmock処理
			// dbHandler
			dbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `users` WHERE twitter_id = ?")).
				WithArgs(args.TwitterID).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "account_name", "twitter_id"}).AddRow(args.UserID, args.TwitterName, args.TwitterAccountName, args.TwitterID))

			// gomock処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// contextHandler
			contextHandler := mock_handler.NewMockContextHandler(mockCtrl)

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)
			loggerHandler.EXPECT().Debug(gomock.Any()).AnyTimes().Return()
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes().Return()

			// twtterHandler
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().SetContextHandler(contextHandler).Return()

			// repository
			authRepository := gateway.NewAuthRepository(
				contextHandler,
				loggerHandler,
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			got, err := authRepository.FindUserByTwitterID(args.TwitterID)

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("FindUserByTwitterID(): %v", err)
			}

			if tt.err != err {
				t.Errorf("FindUserByTwitterID() err = %v, want %v", err, tt.err)
			}

			if *got != *tt.want {
				t.Errorf(
					"FindUserByTwitterID() = %v , want = %v", got, tt.want,
				)
			}

		})
	}
}

// UpsertUserに対するテスト
func TestUpsertUser(t *testing.T) {
	type args struct {
		UserID             uint
		TwitterName        string
		TwitterAccountName string
		TwitterID          string
	}
	table := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "success",
			args: args{
				UserID:             1,
				TwitterName:        "test",
				TwitterAccountName: "Test",
				TwitterID:          "1234567890",
			},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			dbUserHandler, dbMock, err := newMockGormDBUserHandler()

			if err != nil {
				t.Error("sqlmock not work")
			}

			// モックの生成
			// sqlmock処理

			// dbHandler
			dbMock.ExpectBegin()
			dbMock.ExpectExec(
				regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`account_name`,`twitter_id`,`id`) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `updated_at`=?,`deleted_at`=VALUES(`deleted_at`),`name`=VALUES(`name`),`account_name`=VALUES(`account_name`),`twitter_id`=VALUES(`twitter_id`)")).
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), args.TwitterName, args.TwitterAccountName, args.TwitterID, sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))
			dbMock.ExpectCommit()

			// gomock処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// contextHandler
			contextHandler := mock_handler.NewMockContextHandler(mockCtrl)

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes().Return()

			// twtterHandler
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().SetContextHandler(contextHandler).Return()

			// repository
			authRepository := gateway.NewAuthRepository(
				contextHandler,
				loggerHandler,
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			user := entity.NewBlankUser().Update(
				args.UserID,
				args.TwitterName,
				args.TwitterAccountName,
				args.TwitterID,
			)

			err = authRepository.UpsertUser(user)

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("FindUserByTwitterID(): %v", err)
			}

			if tt.err != err {
				t.Errorf("FindUserByTwitterID() err = %v, want %v", err, tt.err)
			}

		})
	}
}

func TestUpdateTwitterApi(t *testing.T) {
	type args struct {
		OAuthToken       string
		OAuthTokenSecret string
	}
	table := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "success",
			args: args{
				OAuthToken:       "token",
				OAuthTokenSecret: "secret",
			},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			dbUserHandler, _, err := newMockGormDBUserHandler()

			if err != nil {
				t.Error("sqlmock not work")
			}

			// gomock処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// contextHandler
			contextHandler := mock_handler.NewMockContextHandler(mockCtrl)

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)

			// twtterHandler
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)
			mockTwitterHandler.EXPECT().UpdateTwitterApi(args.OAuthToken, args.OAuthTokenSecret).Return().AnyTimes()

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().SetContextHandler(contextHandler).Return()

			// repository
			authRepository := gateway.NewAuthRepository(
				contextHandler,
				loggerHandler,
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			// エラーなしで実行できること。
			authRepository.UpdateTwitterApi(args.OAuthToken, args.OAuthTokenSecret)
		})
	}
}

// UpdateSessionに対するテスト
func TestUpdateSession(t *testing.T) {
	type args struct {
		TwitterName        string
		TwitterAccountName string
		TwitterID          string
		UserID             int
		OAuthToken         string
		OAuthTokenSecret   string
	}
	table := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "success",
			args: args{
				TwitterName:        "test",
				TwitterAccountName: "Test",
				TwitterID:          "1234567890",
				UserID:             1,
				OAuthToken:         "token",
				OAuthTokenSecret:   "secret",
			},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			dbUserHandler, _, err := newMockGormDBUserHandler()

			if err != nil {
				t.Error("sqlmock not work")
			}

			// モックの生成
			// gomock処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// contextHandler
			contextHandler := mock_handler.NewMockContextHandler(mockCtrl)

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)
			loggerHandler.EXPECT().Debug(gomock.Any()).AnyTimes().Return()
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes().Return()

			// twtterHandler
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().SetContextHandler(contextHandler).Return()
			sessionHandler.EXPECT().Set(gomock.Any(), gomock.Any()).Return().AnyTimes()
			sessionHandler.EXPECT().Set(gomock.Any(), gomock.Any()).Return().AnyTimes()
			sessionHandler.EXPECT().Save().Return(nil).AnyTimes()

			// repository
			authRepository := gateway.NewAuthRepository(
				contextHandler,
				loggerHandler,
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			err = authRepository.UpdateSession(
				args.OAuthToken,
				args.OAuthTokenSecret,
				args.UserID,
				args.TwitterID,
			)

			if tt.err != err {
				t.Errorf("AuthByCallbackParams() err = %v, want %v", err, tt.err)
			}
		})
	}
}

// Logoutに対するテスト
func TestLogout(t *testing.T) {
	table := []struct {
		name string
		want entity.Auth
		err  error
	}{
		{
			name: "success",
			want: *entity.NewAuth("").SuccessLogout(),
			err:  nil,
		},
	}

	dbUserHandler, _, err := newMockGormDBUserHandler()

	if err != nil {
		t.Error("sqlmock not work")
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// contextHandler
			contextHandler := mock_handler.NewMockContextHandler(mockCtrl)

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)
			loggerHandler.EXPECT().Debug(gomock.Any()).AnyTimes().Return()
			loggerHandler.EXPECT().Debugf(gomock.Any()).AnyTimes().Return()

			// twtterHandler
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().SetContextHandler(contextHandler).Return()
			sessionHandler.EXPECT().Clear().Return()
			sessionHandler.EXPECT().Save().Return(nil)

			authRepository := gateway.NewAuthRepository(
				contextHandler,
				loggerHandler,
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			err := authRepository.Logout()

			if tt.err != err {
				t.Errorf("Logout() err = %v, want %v", err, tt.err)
			}
		})
	}
}
