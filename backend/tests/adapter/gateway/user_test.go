package gateway_test

import (
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

// GetUserByIDに対するテスト
func TestGetUserByID(t *testing.T) {
	type args struct {
		SearchID      int64
		UserID        int64
		UserName      string
		UserTwitterID string
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
				SearchID:      1,
				UserID:        1,
				UserName:      "test",
				UserTwitterID: "1234567890",
			},
			want: &entity.User{
				ID:        1,
				Name:      "test",
				TwitterID: "1234567890",
			},
			err: nil,
		},
		{
			name: "error",
			args: args{
				SearchID:      2,
				UserID:        1,
				UserName:      "test",
				UserTwitterID: "1234567890",
			},
			want: &entity.User{},
			err:  errors.New("user is not found"),
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
				regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
				WithArgs(strconv.FormatInt(args.SearchID, 10)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "twitter_id"}).AddRow(args.UserID, args.UserName, args.UserTwitterID))

			// repository
			userRepository := gateway.NewUserRepository(
				dbUserHandler,
			)

			got, err := userRepository.GetUserByID(strconv.FormatInt(args.SearchID, 10))

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("GetUserByID(): %v", err)
			}

			if err != nil && errors.Is(tt.err, err) {
				t.Errorf("GetUserByID() err = %v, want = %v", err, tt.err)
			}

			if got.ID != tt.want.ID {
				t.Errorf("GetUserByID(%v) user = %v, want = %v", strconv.FormatInt(args.SearchID, 10), got, tt.want)
			}
		})
	}
}
