package database_test

import (
	"errors"
	"os"
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/migration"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
)

var UserDbHandler handler.UserDbHandler

func TestMain(m *testing.M) {
	// 前処理
	config := config.NewConfig(".env.test")

	dbHandler := database.NewGormDbHandler(config)

	UserDbHandler = database.NewUserDbHandler(dbHandler)

	// userテーブルから何も得られない場合はseederを実行
	user := &entity.User{}
	if err := UserDbHandler.Find(user, "id", "1"); err != nil || user.ID == 0 {
		migration.Seeder()
	}

	// トランザクション
	UserDbHandler.Begin()

	status := m.Run()

	// 後処理
	// トランザクションを戻す
	// MySQLではauto_incrementは戻らない
	UserDbHandler.Rollback()

	os.Exit(status)
}

func TestFirst(t *testing.T) {
	type arg struct {
		value string
	}

	table := []struct {
		name     string
		arg      arg
		wantUser entity.User
	}{
		{
			"success",
			arg{
				value: "1",
			},
			entity.User{
				Name:      "name1",
				TwitterID: "1234567890",
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			user := entity.User{}
			UserDbHandler.First(&user, tt.arg.value)

			// 中身の比較
			if user.Name != tt.wantUser.Name || user.TwitterID != tt.wantUser.TwitterID {
				t.Errorf("First()  = %v, want %v", user, tt.wantUser)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type arg struct {
		createSource entity.User
	}

	table := []struct {
		name     string
		arg      arg
		wantUser entity.User
	}{
		{
			"success",
			arg{
				createSource: entity.User{
					Name:        "name2",
					AccountName: "Name2",
					TwitterID:   "1234567892",
				},
			},
			entity.User{
				Name:        "name2",
				AccountName: "Name2",
				TwitterID:   "1234567892",
			},
		},
	}

	UserDbHandler.Transaction(func() error {
		for _, tt := range table {
			t.Run(tt.name, func(t *testing.T) {
				// 作成する
				user := tt.arg.createSource
				UserDbHandler.Create(&user)

				// 中身の比較
				if user.Name != tt.wantUser.Name || user.AccountName != tt.wantUser.AccountName || user.TwitterID != tt.wantUser.TwitterID {
					t.Errorf("Create()  = %v, want %v", user, tt.wantUser)
				}
			})
		}
		return errors.New("rollback")
	})
}

func TestUpdate(t *testing.T) {
	type arg struct {
		column       string
		value        string
		createSource entity.User
	}

	table := []struct {
		name     string
		arg      arg
		wantUser entity.User
	}{
		{
			"success",
			arg{
				column: "id",
				value:  "1",
				createSource: entity.User{
					Name:        "newname1",
					AccountName: "newName1",
					TwitterID:   "1234567890",
				},
			},
			entity.User{
				Name:        "newname1",
				AccountName: "newName1",
				TwitterID:   "1234567890",
			},
		},
	}

	UserDbHandler.Transaction(func() error {
		for _, tt := range table {
			t.Run(tt.name, func(t *testing.T) {
				// 作成する
				user := tt.arg.createSource
				UserDbHandler.Update(&user, tt.arg.value)

				// 中身の比較
				if user.Name != tt.wantUser.Name || user.AccountName != tt.wantUser.AccountName || user.TwitterID != tt.wantUser.TwitterID {
					t.Errorf("Update()  = %v, want %v", user, tt.wantUser)
				}
			})
		}
		return errors.New("rollback")
	})
}

func TestUpsert(t *testing.T) {
	type arg struct {
		column       string
		value        string
		createSource entity.User
	}

	table := []struct {
		name     string
		arg      arg
		wantUser entity.User
	}{
		{
			"success_insert",
			arg{
				column: "id",
				value:  "2",
				createSource: entity.User{
					Name:        "name2",
					AccountName: "Name2",
					TwitterID:   "1234567892",
				},
			},
			entity.User{
				Name:        "name2",
				AccountName: "Name2",
				TwitterID:   "1234567892",
			},
		},
		{
			"success_update",
			arg{
				column: "id",
				value:  "1",
				createSource: entity.User{
					Name:        "newname1",
					AccountName: "newName1",
					TwitterID:   "1234567890",
				},
			},
			entity.User{
				Name:        "newname1",
				AccountName: "newName1",
				TwitterID:   "1234567890",
			},
		},
	}

	UserDbHandler.Transaction(func() error {
		for _, tt := range table {
			t.Run(tt.name, func(t *testing.T) {
				// 作成・更新する
				user := tt.arg.createSource
				UserDbHandler.Upsert(&user, tt.arg.column, tt.arg.value)

				// 中身の比較
				if user.Name != tt.wantUser.Name || user.AccountName != tt.wantUser.AccountName || user.TwitterID != tt.wantUser.TwitterID {
					t.Errorf("Upsert()  = %v, want %v", user, tt.wantUser)
				}
			})
		}

		return errors.New("rollback")
	})
}

func TestFind(t *testing.T) {
	type arg struct {
		column string
		value  string
	}

	table := []struct {
		name     string
		arg      arg
		wantUser entity.User
	}{
		{
			"success",
			arg{
				column: "id",
				value:  "1",
			},
			entity.User{
				Name:      "name1",
				TwitterID: "1234567890",
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			user := entity.User{}
			if err := UserDbHandler.Find(&user, tt.arg.column, tt.arg.value); err != nil {
				t.Errorf("Find(%v, %v) err = %v", tt.arg.column, tt.arg.value, err)
			}

			// 中身の比較
			if user.Name != tt.wantUser.Name || user.TwitterID != tt.wantUser.TwitterID {
				t.Errorf("Find(%v, %v)  = %v, want %v", tt.arg.column, tt.arg.value, user, tt.wantUser)
			}
		})
	}
}

func TestFindByTwitterID(t *testing.T) {
	type arg struct {
		value string
	}

	table := []struct {
		name     string
		arg      arg
		wantUser entity.User
	}{
		{
			"success",
			arg{
				value: "1234567890",
			},
			entity.User{
				Name:      "name1",
				TwitterID: "1234567890",
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			user := entity.User{}
			UserDbHandler.FindByTwitterID(&user, tt.arg.value)

			// 中身の比較
			if user.Name != tt.wantUser.Name || user.TwitterID != tt.wantUser.TwitterID {
				t.Errorf("FindByTwitterID(%v)  = %v, want %v", tt.arg.value, user, tt.wantUser)
			}
		})
	}
}
