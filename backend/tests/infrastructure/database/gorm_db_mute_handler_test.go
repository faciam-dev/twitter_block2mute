package database_test

import (
	"errors"
	"testing"
	"time"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
)

func TestUserMuteFindAllByTwitterID(t *testing.T) {
	type arg struct {
		value string
	}
	nowTime := time.Now()
	table := []struct {
		name         string
		arg          arg
		wantEntities []entity.Mute
	}{
		{
			"success",
			arg{
				value: "1",
			},
			[]entity.Mute{
				*entity.NewMute(1, "12345678904", 1, nowTime, nowTime),
				*entity.NewMute(1, "12345678905", 0, nowTime, nowTime),
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			mutes := []entity.Mute{}
			if err := MuteDBHandler.FindAllByUserID(&mutes, tt.arg.value); err != nil {
				t.Error(err)
			}

			// 中身の比較
			for i, wantMute := range tt.wantEntities {
				isError := true
				for _, gotMute := range mutes {
					if gotMute.GetUserID() == wantMute.GetUserID() &&
						gotMute.GetTargetTwitterID() == wantMute.GetTargetTwitterID() &&
						gotMute.GetFlag() == wantMute.GetFlag() {
						isError = false
						break
					}
				}
				if isError {
					t.Errorf("UserMuteFindAllByTwitterID(%v) got = %v, want = %v[%v]", tt.arg.value, wantMute, mutes, i)
				}
			}
		})
	}
}

func TestUserMuteCreateNew(t *testing.T) {
	type arg struct {
		column       string
		value        string
		createSource entity.Mute
	}
	nowTime := time.Now()
	table := []struct {
		name     string
		arg      arg
		wantMute entity.Mute
	}{
		{
			"success_insert",
			arg{
				column:       "id",
				value:        "3",
				createSource: *entity.NewMute(1, "12345678905", 1, nowTime, nowTime),
			},
			*entity.NewMute(1, "12345678905", 1, nowTime, nowTime),
		},
		{
			"success_update",
			arg{
				column:       "id",
				value:        "4",
				createSource: *entity.NewMute(1, "12345678906", 0, nowTime, nowTime),
			},
			*entity.NewMute(1, "12345678906", 0, nowTime, nowTime),
		},
	}

	DBHandler.Transaction(func(conn handler.DBConnection) error {
		muteDBHandler := database.NewMuteHandler(conn)
		for _, tt := range table {
			t.Run(tt.name, func(t *testing.T) {
				// 作成・更新する
				mutes := []entity.Mute{tt.arg.createSource}
				err := muteDBHandler.CreateNew(&mutes, tt.arg.column, tt.arg.value)
				if err != nil {
					t.Error(err)
				}

				// 中身の比較
				if mutes[0].GetFlag() != tt.wantMute.GetFlag() || mutes[0].GetUserID() != tt.wantMute.GetUserID() || mutes[0].GetTargetTwitterID() != tt.wantMute.GetTargetTwitterID() {
					t.Errorf("UserMuteCreate()  = %v, want %v", mutes[0], tt.wantMute)
				}
			})
		}

		return errors.New("rollback")
	})
	/*
		MuteDBHandler.Transaction(func() error {
			for _, tt := range table {
				t.Run(tt.name, func(t *testing.T) {
					// 作成・更新する
					mutes := []entity.Mute{tt.arg.createSource}
					err := MuteDBHandler.CreateNew(&mutes, tt.arg.column, tt.arg.value)
					if err != nil {
						t.Error(err)
					}

					// 中身の比較
					if mutes[0].GetFlag() != tt.wantMute.GetFlag() || mutes[0].GetUserID() != tt.wantMute.GetUserID() || mutes[0].GetTargetTwitterID() != tt.wantMute.GetTargetTwitterID() {
						t.Errorf("UserMuteCreate()  = %v, want %v", mutes[0], tt.wantMute)
					}
				})
			}

			return errors.New("rollback")
		})
	*/
}
