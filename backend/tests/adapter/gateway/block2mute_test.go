package gateway_test

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/tests/adapter/gateway/mock_handler"
	"github.com/golang/mock/gomock"
)

// Block2MuteGetUserに対するテスト
func TestBlock2MuteGetUser(t *testing.T) {
	type args struct {
		UserID        string
		UserName      string
		UserTwitterID string
		Total         int
		IDs           []string
	}
	table := []struct {
		name     string
		args     args
		wantUser *entity.User
		err      error
	}{
		{
			name: "success_blocked",
			args: args{
				UserID:        "1",
				UserName:      "test",
				UserTwitterID: "1234567890",
				Total:         1,
				IDs: []string{
					"1234567892",
				},
			},
			wantUser: entity.NewBlankUser().Update(1, "test", "test", "1234567890"),
			err:      nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			dbHandler, dbMock, err := newMockGormDbHandler()
			if err != nil {
				t.Error("sqlmock(DB) not work")
			}

			// sqlmock処理
			// userdbHandler
			dbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
				WithArgs(args.UserID).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "account_name", "twitter_id"}).AddRow(args.UserID, args.UserName, args.UserName, args.UserTwitterID))

			// blockDbHandler

			// モックの生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)

			// repository
			repository := gateway.NewBlock2MuteRepository(
				loggerHandler,
				dbHandler,
				mockTwitterHandler,
				sessionHandler,
			)

			got := repository.GetUser(args.UserID)

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Block2MuteGetUser() userDB: %v", err)
			}

			if err != nil && errors.Is(tt.err, err) {
				t.Errorf("Block2MuteGetUser() err = %v, want = %v", err, tt.err)
			}

			if got.GetID() != tt.wantUser.GetID() ||
				got.GetTwitterID() != tt.wantUser.GetTwitterID() ||
				got.GetName() != tt.wantUser.GetName() {
				t.Errorf("Block2MuteGetUser(%v) total = %v, want = %v", args.UserID, got, tt.wantUser)
			}
		})
	}
}

// AuthTwitterに対するテスト
func TestAuthTwitter(t *testing.T) {
	type blocksRow struct {
		ID              string
		UserID          int64
		TargetTwitterID string
		Flag            int
	}

	type mutesRow struct {
		ID              string
		UserID          int64
		TargetTwitterID string
		Flag            int
	}

	type args struct {
		UserID           string
		UserName         string
		UserTwitterID    string
		OAuthToken       string
		OAuthTokenSecret string
		Total            int
		blocks           []blocksRow
		mutes            []mutesRow
	}
	table := []struct {
		name string
		args args
		want *entity.Block2Mute
		err  error
	}{
		{
			name: "success_muted",
			args: args{
				UserID:           "1",
				UserName:         "test",
				UserTwitterID:    "1234567890",
				OAuthToken:       "token",
				OAuthTokenSecret: "secret",
				Total:            1,
				blocks: []blocksRow{
					{
						ID:              "1",
						UserID:          1,
						TargetTwitterID: "1234567891",
						Flag:            0,
					},
				},
				mutes: []mutesRow{
					{
						ID:              "1",
						UserID:          1,
						TargetTwitterID: "1234567892",
						Flag:            0,
					},
				},
			},
			want: entity.NewBlock2Mute(
				1, []string{"1234567891"},
			),
			err: nil,
		},
		{
			name: "success_not_muted",
			args: args{
				UserID:           "1",
				UserName:         "test",
				UserTwitterID:    "1234567890",
				OAuthToken:       "token",
				OAuthTokenSecret: "secret",
				Total:            1,
				blocks:           []blocksRow{},
				mutes:            []mutesRow{},
			},
			want: entity.NewBlock2Mute(
				0, []string{},
			),
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			dbHandler, _, err := newMockGormDbHandler()
			if err != nil {
				t.Error("sqlmock(DB) not work")
			}

			// モックの生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return()
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return()

			// twitterHandler
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)
			mockTwitterHandler.EXPECT().UpdateTwitterApi(args.OAuthToken, args.OAuthTokenSecret).Return().AnyTimes()
			for _, v := range args.blocks {
				mockTwitterHandler.EXPECT().DestroyBlock(args.UserTwitterID, v.TargetTwitterID).Return(nil).AnyTimes()
				mockTwitterHandler.EXPECT().CreateMute(args.UserTwitterID, v.TargetTwitterID).Return(nil).AnyTimes()
			}

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().Get("token").Return(args.OAuthToken).AnyTimes()
			sessionHandler.EXPECT().Get("secret").Return(args.OAuthTokenSecret).AnyTimes()
			sessionHandler.EXPECT().Get("user_id").Return(args.UserID).AnyTimes()

			// repository
			repository := gateway.NewBlock2MuteRepository(
				loggerHandler,
				dbHandler,
				mockTwitterHandler,
				sessionHandler,
			)

			err = repository.AuthTwitter()

			if err != nil && errors.Is(tt.err, err) {
				t.Errorf("All() err = %v, want = %v", err, tt.err)
			}

		})
	}
}

// Allに対するテスト
func TestAll(t *testing.T) {
	type blocksRow struct {
		ID              string
		UserID          int64
		TargetTwitterID string
		Flag            int
	}

	type mutesRow struct {
		ID              string
		UserID          int64
		TargetTwitterID string
		Flag            int
	}

	type args struct {
		UserID           string
		UserName         string
		UserTwitterID    string
		OAuthToken       string
		OAuthTokenSecret string
		Total            int
		blocks           []blocksRow
		mutes            []mutesRow
	}
	table := []struct {
		name string
		args args
		want *entity.Block2Mute
		err  error
	}{
		{
			name: "success_muted",
			args: args{
				UserID:           "1",
				UserName:         "test",
				UserTwitterID:    "1234567890",
				OAuthToken:       "token",
				OAuthTokenSecret: "secret",
				Total:            1,
				blocks: []blocksRow{
					{
						ID:              "1",
						UserID:          1,
						TargetTwitterID: "1234567891",
						Flag:            0,
					},
				},
				mutes: []mutesRow{
					{
						ID:              "1",
						UserID:          1,
						TargetTwitterID: "1234567892",
						Flag:            0,
					},
				},
			},
			want: entity.NewBlock2Mute(
				1, []string{"1234567891"},
			),
			err: nil,
		},
		{
			name: "success_not_muted",
			args: args{
				UserID:           "1",
				UserName:         "test",
				UserTwitterID:    "1234567890",
				OAuthToken:       "token",
				OAuthTokenSecret: "secret",
				Total:            1,
				blocks:           []blocksRow{},
				mutes:            []mutesRow{},
			},
			want: entity.NewBlock2Mute(
				0, []string{},
			),
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			dbHandler, dbMock, err := newMockGormDbHandler()
			if err != nil {
				t.Error("sqlmock(DB) not work")
			}

			// sqlmock処理
			// blockDbHandlerの事前処理
			mockedUserBlocksRow := sqlmock.NewRows([]string{"id", "user_id", "target_twitter_id", "flag"})
			for _, v := range args.blocks {
				mockedUserBlocksRow.AddRow(v.ID, v.UserID, v.TargetTwitterID, v.Flag)
			}

			// 現在のblockを取得
			//blockDbMock.ExpectBegin()
			dbMock.ExpectBegin()
			dbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `user_blocks` WHERE user_id = ? AND `user_blocks`.`deleted_at` IS NULL")).
				WithArgs(args.UserID).
				WillReturnRows(mockedUserBlocksRow)

			// muteDbHandlerの事前処理
			mockedUserMutesRow := sqlmock.NewRows([]string{"id", "user_id", "target_twitter_id", "flag"})
			for _, v := range args.mutes {
				mockedUserMutesRow.AddRow(v.ID, v.UserID, v.TargetTwitterID, v.Flag)
			}

			// 現在のmuteを取得
			dbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `user_mutes` WHERE user_id = ? AND `user_mutes`.`deleted_at` IS NULL")).
				WithArgs(args.UserID).
				WillReturnRows(mockedUserMutesRow)

			// Upsert処理
			// block
			if len(args.blocks) > 0 {
				for _, v := range args.blocks {
					dbMock.ExpectExec(
						regexp.QuoteMeta("INSERT INTO `user_blocks` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`target_twitter_id`,`flag`,`id`) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `updated_at`=?,`deleted_at`=VALUES(`deleted_at`),`user_id`=VALUES(`user_id`),`target_twitter_id`=VALUES(`target_twitter_id`),`flag`=VALUES(`flag`)")).
						WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), v.UserID, v.TargetTwitterID, 1, sqlmock.AnyArg(), sqlmock.AnyArg()).
						WillReturnResult(sqlmock.NewResult(1, 0))
				}
			}

			// mute
			if len(args.mutes) > 0 {
				for _, v := range args.blocks {
					dbMock.ExpectExec(
						regexp.QuoteMeta("INSERT INTO `user_mutes` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`target_twitter_id`,`flag`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `updated_at`=?,`deleted_at`=VALUES(`deleted_at`),`user_id`=VALUES(`user_id`),`target_twitter_id`=VALUES(`target_twitter_id`),`flag`=VALUES(`flag`)")).
						WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), v.UserID, v.TargetTwitterID, 1, sqlmock.AnyArg()).
						WillReturnResult(sqlmock.NewResult(1, 0))
				}
			}

			dbMock.ExpectCommit()

			// モックの生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return()
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return()

			// twitterHandler
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)
			mockTwitterHandler.EXPECT().UpdateTwitterApi(args.OAuthToken, args.OAuthTokenSecret).Return().AnyTimes()
			for _, v := range args.blocks {
				mockTwitterHandler.EXPECT().DestroyBlock(args.UserTwitterID, v.TargetTwitterID).Return(nil).AnyTimes()
				mockTwitterHandler.EXPECT().CreateMute(args.UserTwitterID, v.TargetTwitterID).Return(nil).AnyTimes()
			}

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)

			// repository
			repository := gateway.NewBlock2MuteRepository(
				loggerHandler,
				dbHandler,
				mockTwitterHandler,
				sessionHandler,
			)

			convertedUserID, _ := strconv.Atoi(args.UserID)
			user := entity.NewBlankUser().Update(
				uint(convertedUserID), args.UserName, args.UserName, args.UserTwitterID,
			)

			gotBlock2Mute, err := repository.All(user)

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("All() DB: %v", err)
			}

			if err != nil && errors.Is(tt.err, err) {
				t.Errorf("All() err = %v, want = %v", err, tt.err)
			}

			if !reflect.DeepEqual(gotBlock2Mute.GetSuccessTwitterIDs(), tt.want.GetSuccessTwitterIDs()) ||
				gotBlock2Mute.GetNumberOfSuccess() != tt.want.GetNumberOfSuccess() {
				t.Errorf("All(%v) got = %v, want = %v", args.UserID, gotBlock2Mute, tt.want)
			}

		})
	}
}
