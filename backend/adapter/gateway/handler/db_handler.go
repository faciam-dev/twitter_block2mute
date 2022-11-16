package handler

// 接続
type DbConnection interface {
	GetConnection() interface{}
}

// DB接続に関する操作をするハンドラ
type DBHandler interface {
	Transaction(fn func(DbConnection) error) error
	Begin() DbConnection
	Commit() DbConnection
	Rollback() DbConnection
	Connect() DbConnection
}

// このハンドラは各gatewayで共通の基本的な操作のみを扱う
type CommonDbHandler interface {
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
type UserDbHandler interface {
	CommonDbHandler
	FindByTwitterID(interface{}, string) error
}

// block独自の処理を追加したハンドラ
type BlockDbHandler interface {
	CommonDbHandler
	FindAllByUserID(interface{}, string) error
	CreateNewBlocks(interface{}, string, string) error
}

// mute独自の処理を追加したハンドラ
type MuteDbHandler interface {
	CommonDbHandler
	FindAllByUserID(interface{}, string) error
	CreateNew(interface{}, string, string) error
}
