package database

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"gorm.io/gorm"
)

type GormDBHandler struct {
	Conn *gorm.DB
}

// DBハンドラの初期化
func NewGormDBHandler(gormDBConn handler.DBConnection) handler.DBHandler {
	gormHandler := GormDBHandler{}
	gormHandler.Conn = gormDBConn.GetConnection().(*gorm.DB)

	return gormHandler
}

// transaction
// Begin begins a transaction
func (g GormDBHandler) Begin() handler.DBConnection {
	return NewGormDBConnection(g.Conn.Begin())
}

func (g GormDBHandler) Rollback() handler.DBConnection {
	return NewGormDBConnection(g.Conn.Rollback())
}

func (g GormDBHandler) Commit() handler.DBConnection {
	return NewGormDBConnection(g.Conn.Commit())
}

func (g GormDBHandler) Transaction(fn func(conn handler.DBConnection) error) error {
	return g.Conn.Transaction(func(tx *gorm.DB) error {
		conn := NewGormDBConnection(tx)
		return fn(conn)
	})
}

func (g GormDBHandler) Connect() handler.DBConnection {
	return NewGormDBConnection(g.Conn)
}
