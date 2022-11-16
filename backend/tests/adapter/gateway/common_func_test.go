package gateway_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// mock化したGormをつかったDBへのハンドラを得る
func newMockGormDBHandler() (handler.DBHandler, sqlmock.Sqlmock, error) {
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

	gormDBHandler := database.NewGormDBHandler(
		database.GormDBHandler{Conn: db}.Connect(),
	)

	return gormDBHandler, mock, err
}

// mock化したGormをつかったUserDbへのハンドラを得る
func newMockGormDBUserHandler() (handler.UserDBHandler, sqlmock.Sqlmock, error) {
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

	gormDBUserHandler := database.NewUserDBHandler(
		database.GormDBHandler{Conn: db}.Connect(),
	)

	return gormDBUserHandler, mock, err
}
