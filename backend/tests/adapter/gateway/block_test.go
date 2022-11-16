package gateway_test

import (
	"errors"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/tests/adapter/gateway/mock_handler"
	"github.com/golang/mock/gomock"
)

// GetUserに対するテスト
func TestGetUser(t *testing.T) {
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

			dbHandler, dbMock, err := newMockGormDBHandler()
			if err != nil {
				t.Error("sqlmock(DB) not work")
			}

			// sqlmock処理
			// userdbHandler
			dbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
				WithArgs(args.UserID).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "account_name", "twitter_id"}).AddRow(args.UserID, args.UserName, args.UserName, args.UserTwitterID))

			// blockDBHandler

			// モックの生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)

			// repository
			repository := gateway.NewBlockRepository(
				loggerHandler,
				dbHandler,
				mockTwitterHandler,
				sessionHandler,
			)

			//gotBlocks, gotTotal, err := repository.GetUser(args.UserID)
			got := repository.GetUser(args.UserID)

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("GetUser() userDB: %v", err)
			}

			if err != nil && errors.Is(tt.err, err) {
				t.Errorf("GetUser() err = %v, want = %v", err, tt.err)
			}

			if got.GetID() != tt.wantUser.GetID() ||
				got.GetTwitterID() != tt.wantUser.GetTwitterID() ||
				got.GetName() != tt.wantUser.GetName() {
				t.Errorf("GetUser(%v) total = %v, want = %v", args.UserID, got, tt.wantUser)
			}
		})
	}
}

// GetGetBlocksに対するテスト
func TestGetBlocks(t *testing.T) {
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
	nowTime := time.Now()
	table := []struct {
		name       string
		args       args
		wantBlocks *entity.Blocks
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
			wantBlocks: &entity.Blocks{
				*entity.NewBlock(1, "1234567892", 0, nowTime, nowTime),
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
			wantBlocks: &entity.Blocks{
				*entity.NewBlock(1, "1234567892", 0, nowTime, nowTime),
			},
			wantTotal: 1,
			err:       nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			dbHandler, _, err := newMockGormDBHandler()
			if err != nil {
				t.Error("sqlmock(DB) not work")
			}

			// モックの生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockTwitterUser := mock_handler.NewMockTwitterUserIds(mockCtrl)
			mockTwitterUser.EXPECT().GetTotal().Return(args.Total).AnyTimes()
			mockTwitterUser.EXPECT().GetTwitterIDs().Return(args.IDs).AnyTimes()

			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)
			mockTwitterHandler.EXPECT().GetBlockedUser(args.UserTwitterID).Return(mockTwitterUser, tt.err).AnyTimes()
			mockTwitterHandler.EXPECT().UpdateTwitterApi(args.OAuthToken, args.OAuthTokenSecret).Return().AnyTimes()

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return()
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return()

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().Get("token").Return(args.OAuthToken).AnyTimes()
			sessionHandler.EXPECT().Get("secret").Return(args.OAuthTokenSecret).AnyTimes()
			sessionHandler.EXPECT().Get("user_id").Return(args.UserID).AnyTimes()

			// repository
			repository := gateway.NewBlockRepository(
				loggerHandler,
				dbHandler,
				mockTwitterHandler,
				sessionHandler,
			)

			convertedUserID, _ := strconv.ParseUint(args.UserID, 10, 0)
			gotBlocks, gotTotal, err := repository.GetBlocks(
				entity.NewBlankUser().Update(uint(convertedUserID), args.UserName, args.UserName, args.UserTwitterID),
			)

			if err != nil && errors.Is(tt.err, err) {
				t.Errorf("GetBlocks() err = %v, want = %v", err, tt.err)
			}

			for i, gotBlock := range *gotBlocks {
				isError := true
				for _, wantBlock := range *tt.wantBlocks {
					if gotBlock.GetUserID() == wantBlock.GetUserID() &&
						gotBlock.GetTargetTwitterID() == wantBlock.GetTargetTwitterID() &&
						gotBlock.GetFlag() == wantBlock.GetFlag() {
						isError = false
						break
					}
				}
				if isError {
					t.Errorf("GetBlocks(%v) got = %v, want = %v[%v]", args.UserID, gotBlock, *tt.wantBlocks, i)
				}
			}

			if gotTotal != tt.wantTotal {
				t.Errorf("GetBlocks(%v) total = %v, want = %v", args.UserID, gotTotal, tt.wantTotal)
			}
		})
	}
}

// TxUpdateAndDeleteBlocksに対するテスト
func TestTxUpdateAndDeleteBlocks(t *testing.T) {
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
	nowTime := time.Now()
	table := []struct {
		name      string
		args      args
		argBlocks *entity.Blocks
		wantTotal int
		err       error
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
						TargetTwitterID: "1234567890",
						Flag:            1,
					},
				},
			},
			argBlocks: &entity.Blocks{
				*entity.NewBlock(1, "1234567892", 0, nowTime, nowTime),
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
			argBlocks: &entity.Blocks{},
			wantTotal: 1,
			err:       nil,
		},
		{
			name: "fail_not_blocked",
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
			argBlocks: &entity.Blocks{
				*entity.NewBlock(1, "1234567892", 0, nowTime, nowTime),
			},
			wantTotal: 1,
			err:       nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			dbHandler, dbMock, err := newMockGormDBHandler()
			if err != nil {
				t.Error("sqlmock(DB) not work")
			}

			// sqlmock処理
			// blockDBHandler
			mockedUserBlocksRow := sqlmock.NewRows([]string{"id", "user_id", "target_twitter_id", "flag"})
			for _, v := range args.blocks {
				mockedUserBlocksRow.AddRow(v.ID, v.UserID, v.TargetTwitterID, v.Flag)
			}
			dbMock.ExpectBegin()
			dbMock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `user_blocks` WHERE user_id = ? AND `user_blocks`.`deleted_at` IS NULL")).
				WithArgs(args.UserID).
				WillReturnRows(mockedUserBlocksRow)

			// 存在しないblockを削除する処理
			if len(args.blocks) > 0 {
				dbMock.ExpectExec(
					regexp.QuoteMeta("UPDATE `user_blocks` SET `deleted_at`=? WHERE `user_blocks`.`id` = ? AND `user_blocks`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			// Upsert処理。ブロックしたものがある場合だけ実行される。
			argBlocks := entity.Blocks{}
			if len(args.IDs) > 0 {
				convertedUserID64, _ := strconv.ParseInt(args.UserID, 10, 64)
				dbMock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO `user_blocks` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`target_twitter_id`,`flag`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `updated_at`=?,`deleted_at`=VALUES(`deleted_at`),`user_id`=VALUES(`user_id`),`target_twitter_id`=VALUES(`target_twitter_id`),`flag`=VALUES(`flag`)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), convertedUserID64, args.IDs[0], 0, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(2, 0))
				argBlocks = append(argBlocks, *entity.NewBlock(1, args.IDs[0], 0, nowTime, nowTime))
			}
			dbMock.ExpectCommit()

			// モックの生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockTwitterUser := mock_handler.NewMockTwitterUserIds(mockCtrl)
			mockTwitterUser.EXPECT().GetTotal().Return(args.Total).AnyTimes()
			mockTwitterUser.EXPECT().GetTwitterIDs().Return(args.IDs).AnyTimes()

			mockTwitterHandler := mock_handler.NewMockTwitterHandler(mockCtrl)
			mockTwitterHandler.EXPECT().GetBlockedUser(args.UserTwitterID).Return(mockTwitterUser, tt.err).AnyTimes()
			mockTwitterHandler.EXPECT().UpdateTwitterApi(args.OAuthToken, args.OAuthTokenSecret).Return().AnyTimes()

			// loggerHandler
			loggerHandler := mock_handler.NewMockLoggerHandler(mockCtrl)
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes().Return()
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return()
			loggerHandler.EXPECT().Debugf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return()
			loggerHandler.EXPECT().Errorw(gomock.Any(), gomock.Any()).AnyTimes().Return()

			// sessionHandler
			sessionHandler := mock_handler.NewMockSessionHandler(mockCtrl)
			sessionHandler.EXPECT().Get("token").Return(args.OAuthToken).AnyTimes()
			sessionHandler.EXPECT().Get("secret").Return(args.OAuthTokenSecret).AnyTimes()
			sessionHandler.EXPECT().Get("user_id").Return(args.UserID).AnyTimes()

			// repository
			repository := gateway.NewBlockRepository(
				loggerHandler,
				dbHandler,
				mockTwitterHandler,
				sessionHandler,
			)

			convertedUserID, _ := strconv.ParseUint(args.UserID, 10, 0)
			err = repository.TxUpdateAndDeleteBlocks(
				entity.NewBlankUser().Update(uint(convertedUserID), args.UserName, args.UserName, args.UserTwitterID), &argBlocks,
			)

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("TxUpdateAndDeleteBlocks() DB: %v", err)
			}

			if err != nil && errors.Is(tt.err, err) {
				t.Errorf("TxUpdateAndDeleteBlocks() err = %v, want = %v", err, tt.err)
			}
		})
	}
}
