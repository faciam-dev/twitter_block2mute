package gateway_test

import (
	"errors"
	"regexp"
	"strconv"
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

// mock化したGormをつかったBlockDbへのハンドラを得る
func newMockGormDbBlockHandler() (handler.BlockDbHandler, sqlmock.Sqlmock, error) {
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

	gormDbBlockHandler := database.NewBlockDbHandler(
		database.GormDbHandler{Conn: db},
	)

	return gormDbBlockHandler, mock, err
}

// GetUserIDsに対するテスト
func TestGetUserIDs(t *testing.T) {
	type blocksRow struct {
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
		IDs              []string
		blocks           []blocksRow
	}
	table := []struct {
		name       string
		args       args
		wantBlocks *[]entity.Block
		wantTotal  int
		err        error
	}{
		{
			name: "success_blocked",
			args: args{
				UserID:           "1",
				UserName:         "test",
				UserTwitterID:    "1234567890",
				OAuthToken:       "token",
				OAuthTokenSecret: "secret",
				Total:            1,
				IDs: []string{
					"1234567892",
				},
				blocks: []blocksRow{
					{
						ID:              "1",
						UserID:          1,
						TargetTwitterID: "1234567891",
						Flag:            1,
					},
				},
			},
			wantBlocks: &[]entity.Block{
				{
					ID:              1,
					UserID:          1,
					TargetTwitterID: "1234567892",
					Flag:            0,
				},
			},
			wantTotal: 1,
			err:       nil,
		},
		{
			name: "success_not_blocked",
			args: args{
				UserID:           "1",
				UserName:         "test",
				UserTwitterID:    "1234567890",
				OAuthToken:       "token",
				OAuthTokenSecret: "secret",
				Total:            1,
				IDs:              []string{},
				blocks:           []blocksRow{},
			},
			wantBlocks: &[]entity.Block{
				{
					ID:              1,
					UserID:          1,
					TargetTwitterID: "1234567892",
					Flag:            0,
				},
			},
			wantTotal: 1,
			err:       nil,
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

			// sqlmock処理
			// userdbHandler
			userDbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
				WithArgs(args.UserID).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "twitter_id"}).AddRow(args.UserID, args.UserName, args.UserTwitterID))

			// blockDbHandler
			mockedUserBlocksRow := sqlmock.NewRows([]string{"id", "user_id", "target_twitter_id", "flag"})
			for _, v := range args.blocks {
				mockedUserBlocksRow.AddRow(v.ID, v.UserID, v.TargetTwitterID, v.Flag)
			}
			//blockDbMock.ExpectQuery("").WithArgs(sqlmock.AnyArg)
			blockDbMock.ExpectBegin()
			blockDbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `user_blocks` WHERE user_id = ? AND `user_blocks`.`deleted_at` IS NULL")).
				WithArgs(args.UserID).
				WillReturnRows(mockedUserBlocksRow)

			// 存在しないblockを削除する処理
			if len(args.blocks) > 0 {
				blockDbMock.ExpectBegin()
				blockDbMock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_blocks` SET `deleted_at`=? WHERE `user_blocks`.`id` = ? AND `user_blocks`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				blockDbMock.ExpectCommit()
			}

			// Upsert処理。ブロックしたものがある場合だけ実行される。
			if len(args.IDs) > 0 {
				convertedUserID64, _ := strconv.ParseInt(args.UserID, 10, 64)
				blockDbMock.ExpectBegin()
				blockDbMock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO `user_blocks` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`target_twitter_id`,`flag`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `updated_at`=?,`deleted_at`=VALUES(`deleted_at`),`user_id`=VALUES(`user_id`),`target_twitter_id`=VALUES(`target_twitter_id`),`flag`=VALUES(`flag`)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), convertedUserID64, args.IDs[0], 0, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(2, 0))
				//blockDbMock.ExpectRollback()
				blockDbMock.ExpectCommit()
			}
			//blockDbMock.ExpectBegin()

			blockDbMock.ExpectCommit()

			// モックの生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockTwitterUser := mock_handler.NewMockTwitterUserIds(mockCtrl)
			mockTwitterUser.EXPECT().GetTotal().Return(args.Total).AnyTimes()
			mockTwitterUser.EXPECT().GetTwitterIDs().Return(args.IDs).AnyTimes()

			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)
			mockTwitterHandler.EXPECT().GetBlockedUser(args.UserTwitterID).Return(mockTwitterUser, tt.err).AnyTimes()
			mockTwitterHandler.EXPECT().UpdateTwitterApi(args.OAuthToken, args.OAuthTokenSecret).Return().AnyTimes()

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().Get("token").Return(args.OAuthToken).AnyTimes()
			sessionHandler.EXPECT().Get("secret").Return(args.OAuthTokenSecret).AnyTimes()
			sessionHandler.EXPECT().Get("user_id").Return(args.UserID).AnyTimes()

			// repository
			repository := gateway.NewBlockRepository(
				blockDbHandler,
				userDbHandler,
				mockTwitterHandler,
				sessionHandler,
			)

			gotBlocks, gotTotal, err := repository.GetUserIDs(args.UserID)

			if err := userDbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("GetUserIDs() userDB: %v", err)
			}

			if err := blockDbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("GetUserIDs() blockDB: %v", err)
			}

			if err != nil && errors.Is(tt.err, err) {
				t.Errorf("GetUserIDs() err = %v, want = %v", err, tt.err)
			}

			for i, gotBlock := range *gotBlocks {
				isError := true
				for _, wantBlock := range *tt.wantBlocks {
					if gotBlock.UserID == wantBlock.UserID &&
						gotBlock.TargetTwitterID == wantBlock.TargetTwitterID &&
						gotBlock.Flag == wantBlock.Flag {
						isError = false
						break
					}
				}
				if isError {
					t.Errorf("GetUserIDs(%v) got = %v, want = %v[%v]", args.UserID, gotBlock, *tt.wantBlocks, i)
				}
			}

			if gotTotal != tt.wantTotal {
				t.Errorf("GetUserIDs(%v) total = %v, want = %v", args.UserID, gotTotal, tt.wantTotal)
			}
		})
	}
}
