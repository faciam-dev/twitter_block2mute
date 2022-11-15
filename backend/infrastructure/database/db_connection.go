package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
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

// DB接続を得る
func (g *GormDbConnection) GetConnection() interface{} {
	return g.db
}
