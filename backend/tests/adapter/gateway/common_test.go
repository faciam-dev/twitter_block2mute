package gateway_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// mock化したGormをつかったDBへのハンドラを得る
func newMockGormDbHandler() (handler.DBConnectionHandler, sqlmock.Sqlmock, error) {
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

	gormDbHandler := database.NewGormDbHandlerByConnection(db)

	return gormDbHandler, mock, err
}

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
		database.GormDbHandler{Conn: db}.Connect(),
	)

	return gormDbUserHandler, mock, err
}
