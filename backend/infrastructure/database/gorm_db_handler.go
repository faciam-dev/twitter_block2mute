package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDbHandler struct {
	DbInfo
	Conn *gorm.DB
}

func NewGormDbHandler(gormDBConn handler.DbConnection) handler.DBHandler {
	gormHandler := GormDbHandler{}
	gormHandler.Conn = gormDBConn.GetConnection().(*gorm.DB)

	return gormHandler
}

func NewGormDbHandler_(config *config.Config) handler.DBHandler {
	return newGormDbHandler(DbInfo{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		Username: config.DB.Username,
		Password: config.DB.Password,
		DBName:   config.DB.DBName,
	})
}

func newGormDbHandler(dbHandler DbInfo) handler.DBHandler {
	dsn := dbHandler.Username + ":" + dbHandler.Password + "@tcp(" + dbHandler.Host + ":" + dbHandler.Port + ")/" + dbHandler.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// https://github.com/go-sql-driver/mysql#examples
	//db, err := gorm.Open("mysql", d.Username + ":" + d.Password + "@tcp(" + d.Host + ")/" + d.DBName + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}

	gormHandler := GormDbHandler{
		DbInfo: dbHandler,
	}

	gormHandler.Conn = db

	return gormHandler
}

// transaction
// Begin begins a transaction
func (g GormDbHandler) Begin() handler.DbConnection {
	return NewGormDbConnection(g.Conn.Begin())
}

func (g GormDbHandler) Rollback() handler.DbConnection {
	return NewGormDbConnection(g.Conn.Rollback())
}

func (g GormDbHandler) Commit() handler.DbConnection {
	return NewGormDbConnection(g.Conn.Commit())
}

func (g GormDbHandler) Transaction(fn func(conn handler.DbConnection) error) error {
	return g.Conn.Transaction(func(tx *gorm.DB) error {
		conn := NewGormDbConnection(tx)
		return fn(conn)
	})
}

func (g GormDbHandler) Connect() handler.DbConnection {
	return NewGormDbConnection(g.Conn)
}
