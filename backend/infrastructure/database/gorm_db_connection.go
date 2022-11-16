package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDbConnection struct {
	db *gorm.DB
}

// DB接続を作成する
func NewGormDbConnection(db *gorm.DB) handler.DbConnection {
	dbConnection := new(GormDbConnection)
	dbConnection.db = db

	return dbConnection
}

// DB接続をConfigから作成する
func NewGormDbConnectionByConfig(config *config.Config) handler.DbConnection {
	dsn := config.DB.Username + ":" + config.DB.Password + "@tcp(" + config.DB.Host + ":" + config.DB.Port + ")/" + config.DB.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// https://github.com/go-sql-driver/mysql#examples
	if err != nil {
		panic(err.Error())
	}

	dbConnection := new(GormDbConnection)
	dbConnection.db = db

	return dbConnection
}

// DB接続を得る
func (g *GormDbConnection) GetConnection() interface{} {
	return g.db
}
