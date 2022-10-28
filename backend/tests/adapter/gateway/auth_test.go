package gateway_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/common"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
	"github.com/faciam_dev/twitter_block2mute/backend/tests/adapter/gateway/mock_handler"
)

// mock化したGormをつかったUserDbへのハンドラを得る
func newMockGormDbUserHandler() (handler.UserDbHandler, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, mock, err
	}

	db, err := gorm.Open(mysql.New(
		mysql.Config{
			DriverName:                "mysql",
			Conn:                      mockDB,
			SkipInitializeWithVersion: true,
		}),
		&gorm.Config{},
	)

	gormDbUserHandler := database.NewUserDbHandler(
		database.GormDbHandler{Conn: db},
	)

	return gormDbUserHandler, mock, err
}

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
			want: entity.Auth{
				Authenticated: 1,
			},
			err: nil,
		},
		{
			name: "success",
			args: args{
				SessionToken:       nil,
				SessionTokenSecret: nil,
				TwitterName:        "test",
				TwitterAccountName: "Test",
				TwitterID:          "1234567890",
				OAuthToken:         "token",
				OAuthTokenSecret:   "secret",
			},
			want: entity.Auth{
				Authenticated: 0,
			},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			// モックの生成
			// sqlmock処理
			dbUserHandler /*dbMock*/, _, err := newMockGormDbUserHandler()

			if err != nil {
				t.Error("sqlmock not work")
			}

			// モック処理
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// contextHandler
			contextHandler := mock_handler.NewMockContextHandler(mockCtrl)

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
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			got, err := authRepository.IsAuthenticated()
			if tt.err != err {
				t.Errorf("IsAuthenticated() err = %v, want %v", err, tt.err)
			}

			if got.Authenticated != tt.want.Authenticated {
				t.Errorf("IsAuthenticated() Authenticated = %v, want %v", got, tt.want)
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
			want: entity.Auth{
				Authenticated: 0,
				AuthUrl:       "https://api.twitter.com/oauth/authenticate",
			},
			err: nil,
		},
	}

	dbUserHandler, _, err := newMockGormDbUserHandler()

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

			// twtterHandler
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)
			mockTwitterHandler.EXPECT().AuthorizationURL().Return(args.AuthUrl+"?token="+common.RandomString(16), nil)

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().SetContextHandler(contextHandler).Return()

			authRepository := gateway.NewAuthRepository(
				contextHandler,
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			got, err := authRepository.Auth()

			if tt.err != err {
				t.Errorf("Auth() err = %v, want %v", err, tt.err)
			}

			if got.Authenticated != tt.want.Authenticated || strings.Index(got.AuthUrl, tt.want.AuthUrl) > 0 {
				t.Errorf("Auth() auth = %v, want %v", got, tt.want)
			}
		})
	}
}

// Callbackに対するテスト
func TestCallback(t *testing.T) {
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
			want: entity.Auth{
				Authenticated: 1,
			},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			dbUserHandler, dbMock, err := newMockGormDbUserHandler()

			if err != nil {
				t.Error("sqlmock not work")
			}

			// モックの生成
			// sqlmock処理
			// dbHandler
			dbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `users` WHERE twitter_id = ?")).
				WithArgs(args.TwitterID).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "account_name", "twitter_id"}).AddRow(1, args.TwitterName, args.TwitterAccountName, args.TwitterID))

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
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			got, err := authRepository.Callback(args.OAuthToken, args.OAuthTokenSecret, args.TwitterID, args.TwitterName)

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("IsAuthenticated(): %v", err)
			}

			if tt.err != err {
				t.Errorf("IsAuthenticated() err = %v, want %v", err, tt.err)
			}

			if got.Authenticated != tt.want.Authenticated {
				t.Errorf("IsAuthenticated() Authenticated = %v, want %v", got, tt.want)
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
			want: entity.Auth{
				Logout: 1,
			},
			err: nil,
		},
	}

	dbUserHandler, _, err := newMockGormDbUserHandler()

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

			// twtterHandler
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().SetContextHandler(contextHandler).Return()
			sessionHandler.EXPECT().Clear().Return()

			authRepository := gateway.NewAuthRepository(
				contextHandler,
				mockTwitterHandler,
				sessionHandler,
				dbUserHandler,
			)

			got, err := authRepository.Logout()

			if tt.err != err {
				t.Errorf("Auth() err = %v, want %v", err, tt.err)
			}

			if got.Logout != tt.want.Logout {
				t.Errorf("Auth() auth = %v, want %v", got, tt.want)
			}
		})
	}
}
