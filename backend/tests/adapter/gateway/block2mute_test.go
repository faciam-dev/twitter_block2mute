package gateway_test

import (
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
	"github.com/faciam_dev/twitter_block2mute/backend/tests/adapter/gateway/mock_handler"
	"github.com/golang/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// mock化したGormをつかったMuteDbへのハンドラを得る
func newMockGormDbMuteHandler() (handler.MuteDbHandler, sqlmock.Sqlmock, error) {
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

	gormDbMuteHandler := database.NewMuteHandler(
		database.GormDbHandler{Conn: db},
	)

	return gormDbMuteHandler, mock, err
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
			want: &entity.Block2Mute{
				NumberOfSuccess:   1,
				SuccessTwitterIDs: []string{"1234567891"},
			},
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
			want: &entity.Block2Mute{
				NumberOfSuccess:   0,
				SuccessTwitterIDs: []string{},
			},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			userDbHandler, userDbMock, err := newMockGormDbUserHandler()
			if err != nil {
				t.Error("sqlmock(User) not work")
			}

			blockDbHandler, blockDbMock, err := newMockGormDbBlockHandler()
			if err != nil {
				t.Error("sqlmock(Block) not work")
			}

			muteDbHandler, muteDbMock, err := newMockGormDbMuteHandler()
			if err != nil {
				t.Error("sqlmock(Mute) not work")
			}

			// sqlmock処理
			// userdbHandler
			userDbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
				WithArgs(args.UserID).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "twitter_id"}).AddRow(args.UserID, args.UserName, args.UserTwitterID))

			// blockDbHandlerの事前処理
			mockedUserBlocksRow := sqlmock.NewRows([]string{"id", "user_id", "target_twitter_id", "flag"})
			for _, v := range args.blocks {
				mockedUserBlocksRow.AddRow(v.ID, v.UserID, v.TargetTwitterID, v.Flag)
			}

			// 現在のblockを取得
			//blockDbMock.ExpectBegin()
			muteDbMock.ExpectBegin()
			blockDbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `user_blocks` WHERE user_id = ? AND `user_blocks`.`deleted_at` IS NULL")).
				WithArgs(args.UserID).
				WillReturnRows(mockedUserBlocksRow)

			// muteDbHandlerの事前処理
			mockedUserMutesRow := sqlmock.NewRows([]string{"id", "user_id", "target_twitter_id", "flag"})
			for _, v := range args.mutes {
				mockedUserMutesRow.AddRow(v.ID, v.UserID, v.TargetTwitterID, v.Flag)
			}

			// 現在のmuteを取得
			muteDbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `user_mutes` WHERE user_id = ? AND `user_mutes`.`deleted_at` IS NULL")).
				WithArgs(args.UserID).
				WillReturnRows(mockedUserMutesRow)

			// Upsert処理
			// block
			if len(args.blocks) > 0 {
				for _, v := range args.blocks {
					//convertedUserID64, _ := strconv.ParseInt(args.UserID, 10, 64)
					blockDbMock.ExpectBegin()
					blockDbMock.ExpectExec(
						regexp.QuoteMeta("INSERT INTO `user_blocks` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`target_twitter_id`,`flag`,`id`) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `updated_at`=?,`deleted_at`=VALUES(`deleted_at`),`user_id`=VALUES(`user_id`),`target_twitter_id`=VALUES(`target_twitter_id`),`flag`=VALUES(`flag`)")).
						WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), v.UserID, v.TargetTwitterID, 1, sqlmock.AnyArg(), sqlmock.AnyArg()).
						WillReturnResult(sqlmock.NewResult(1, 0))
					//muteDbMock.ExpectRollback()
					blockDbMock.ExpectCommit()
				}
			}

			// mute
			if len(args.mutes) > 0 {
				for _, v := range args.blocks {
					//convertedUserID64, _ := strconv.ParseInt(args.UserID, 10, 64)
					muteDbMock.ExpectBegin()
					muteDbMock.ExpectExec(
						regexp.QuoteMeta("INSERT INTO `user_mutes` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`target_twitter_id`,`flag`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `updated_at`=?,`deleted_at`=VALUES(`deleted_at`),`user_id`=VALUES(`user_id`),`target_twitter_id`=VALUES(`target_twitter_id`),`flag`=VALUES(`flag`)")).
						WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), v.UserID, v.TargetTwitterID, 1, sqlmock.AnyArg()).
						WillReturnResult(sqlmock.NewResult(1, 0))
					//muteDbMock.ExpectRollback()
					muteDbMock.ExpectCommit()
				}
			}
			//muteDbMock.ExpectBegin()

			muteDbMock.ExpectCommit()
			//blockDbMock.ExpectCommit()

			// モックの生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

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
				blockDbHandler,
				userDbHandler,
				muteDbHandler,
				mockTwitterHandler,
				sessionHandler,
			)

			gotBlock2Mute, err := repository.All(args.UserID)

			if err := userDbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("All() userDB: %v", err)
			}

			if err := blockDbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("All() blockDB: %v", err)
			}

			if err := muteDbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("All() muteDB: %v", err)
			}

			if err != nil && errors.Is(tt.err, err) {
				t.Errorf("All() err = %v, want = %v", err, tt.err)
			}

			if !reflect.DeepEqual(gotBlock2Mute.SuccessTwitterIDs, tt.want.SuccessTwitterIDs) ||
				gotBlock2Mute.NumberOfSuccess != tt.want.NumberOfSuccess {
				t.Errorf("All(%v) got = %v, want = %v", args.UserID, gotBlock2Mute, tt.want)
			}

		})
	}
}
