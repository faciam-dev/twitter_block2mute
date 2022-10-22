package database_test

import (
	"errors"
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

func TestUserMuteFindAllByTwitterID(t *testing.T) {
	type arg struct {
		value string
	}

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
				{
					UserID:          1,
					TargetTwitterID: "12345678904",
					Flag:            1,
				},
				{
					UserID:          1,
					TargetTwitterID: "12345678905",
					Flag:            0,
				},
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			mutes := []entity.Mute{}
			if err := MuteDbHandler.FindAllByUserID(&mutes, tt.arg.value); err != nil {
				t.Error(err)
			}

			// 中身の比較
			for i, wantMute := range tt.wantEntities {
				isError := true
				for _, gotMute := range mutes {
					if gotMute.UserID == wantMute.UserID &&
						gotMute.TargetTwitterID == wantMute.TargetTwitterID &&
						gotMute.Flag == wantMute.Flag {
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

	table := []struct {
		name     string
		arg      arg
		wantMute entity.Mute
	}{
		{
			"success_insert",
			arg{
				column: "id",
				value:  "3",
				createSource: entity.Mute{
					UserID:          1,
					TargetTwitterID: "12345678905",
					Flag:            1,
				},
			},
			entity.Mute{
				UserID:          1,
				TargetTwitterID: "12345678905",
				Flag:            1,
			},
		},
		{
			"success_update",
			arg{
				column: "id",
				value:  "4",
				createSource: entity.Mute{
					UserID:          1,
					TargetTwitterID: "12345678906",
					Flag:            0,
				},
			},
			entity.Mute{
				UserID:          1,
				TargetTwitterID: "12345678906",
				Flag:            0,
			},
		},
	}

	MuteDbHandler.Transaction(func() error {
		for _, tt := range table {
			t.Run(tt.name, func(t *testing.T) {
				// 作成・更新する
				mutes := []entity.Mute{tt.arg.createSource}
				err := MuteDbHandler.CreateNew(&mutes, tt.arg.column, tt.arg.value)
				if err != nil {
					t.Error(err)
				}

				// 中身の比較
				if mutes[0].Flag != tt.wantMute.Flag || mutes[0].UserID != tt.wantMute.UserID || mutes[0].TargetTwitterID != tt.wantMute.TargetTwitterID {
					t.Errorf("UserMuteCreate()  = %v, want %v", mutes[0], tt.wantMute)
				}
			})
		}

		return errors.New("rollback")
	})
}
