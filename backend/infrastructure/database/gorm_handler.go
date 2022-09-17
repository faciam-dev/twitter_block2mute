package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormHandler struct {
    Host string
    Port string
    Username string
    Password string
    DBName string
    Conn *gorm.DB
}

func NewGormHandler(config *config.Config) GormHandler {
    return newGormHandler(GormHandler{
        Host: config.DB.Host,
        Port: config.DB.Port,
        Username: config.DB.Username,
        Password: config.DB.Password,
        DBName: config.DB.DBName,
    })
}

func newGormHandler(gormHandler GormHandler) GormHandler {
    dsn := gormHandler.Username + ":" + gormHandler.Password + "@tcp(" + gormHandler.Host + ":" + gormHandler.Port + ")/" + gormHandler.DBName + "?parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    // https://github.com/go-sql-driver/mysql#examples
    //db, err := gorm.Open("mysql", d.Username + ":" + d.Password + "@tcp(" + d.Host + ")/" + d.DBName + "?charset=utf8&parseTime=True&loc=Local")
    if err != nil {
        panic(err.Error())
    }

    gormHandler.Conn = db

    return gormHandler
}

// Begin begins a transaction
func (g *GormHandler) Begin() *gorm.DB {
    return g.Conn.Begin()
}

func (g *GormHandler) Connect() *gorm.DB {
    return g.Conn
}
