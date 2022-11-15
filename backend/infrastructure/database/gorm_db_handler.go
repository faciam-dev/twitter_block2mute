package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDbHandler struct {
	DbHandler
	Conn *gorm.DB
}

func NewGormDbHandler(config *config.Config) *GormDbHandler {
	return newGormDbHandler(DbHandler{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		Username: config.DB.Username,
		Password: config.DB.Password,
		DBName:   config.DB.DBName,
	})
}

func newGormDbHandler(dbHandler DbHandler) *GormDbHandler {
	dsn := dbHandler.Username + ":" + dbHandler.Password + "@tcp(" + dbHandler.Host + ":" + dbHandler.Port + ")/" + dbHandler.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// https://github.com/go-sql-driver/mysql#examples
	//db, err := gorm.Open("mysql", d.Username + ":" + d.Password + "@tcp(" + d.Host + ")/" + d.DBName + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}

	gormHandler := GormDbHandler{
		DbHandler: dbHandler,
	}

	gormHandler.Conn = db

	return &gormHandler
}

// transaction
// Begin begins a transaction
func (g *GormDbHandler) Begin() *gorm.DB {
	return g.Conn.Begin()
}

func (g *GormDbHandler) Rollback() *gorm.DB {
	return g.Conn.Rollback()
}

func (g *GormDbHandler) Commit() *gorm.DB {
	return g.Conn.Commit()
}

func (g *GormDbHandler) Transaction(fn func() error) error {
	return g.Conn.Transaction(func(tx *gorm.DB) error {
		return fn()
	})
}

func (g *GormDbHandler) Connect() *gorm.DB {
	return g.Conn
}
