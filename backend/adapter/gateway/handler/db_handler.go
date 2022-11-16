package handler

// 接続
type DBConnection interface {
	GetConnection() interface{}
}

// DB接続に関する操作をするハンドラ
type DBHandler interface {
	Transaction(fn func(DBConnection) error) error
	Begin() DBConnection
	Commit() DBConnection
	Rollback() DBConnection
	Connect() DBConnection
}

// このハンドラは各gatewayで共通の基本的な操作のみを扱う
type CommonDBHandler interface {
	First(interface{}, string) error
	Create(interface{}) error
	Find(interface{}, string, string) error
	FindAll(interface{}, string, string) error
	Update(interface{}, string) error
	Upsert(interface{}, string, string) error
	Delete(interface{}, string) error
	DeleteByIds(interface{}, []uint) error
	DeleteAll(interface{}, string, string) error
}

// user独自の処理を追加したハンドラ
type UserDBHandler interface {
	CommonDBHandler
	FindByTwitterID(interface{}, string) error
}

// block独自の処理を追加したハンドラ
type BlockDBHandler interface {
	CommonDBHandler
	FindAllByUserID(interface{}, string) error
	CreateNewBlocks(interface{}, string, string) error
}

// mute独自の処理を追加したハンドラ
type MuteDBHandler interface {
	CommonDBHandler
	FindAllByUserID(interface{}, string) error
	CreateNew(interface{}, string, string) error
}
