package handler

// このハンドラは各gatewayで共通の基本的な操作のみを扱う
type DbHandler interface {
	Transaction(fn func() error) error
	Begin()
	Commit()
	Rollback()
	First(interface{}, string) error
	Create(interface{}) error
	Update(interface{}, string) error
	Upsert(interface{}, string, string) error
	Find(interface{}, string, string) error
	FindAll(interface{}, string, string) error
}

// user独自の処理を追加したハンドラ
type UserDbHandler interface {
	DbHandler
	FindByTwitterID(interface{}, string) error
}

// block独自の処理を追加したハンドラ
type BlockDbHandler interface {
	DbHandler
	FindAllByUserID(interface{}, string) error
	CreateNewBlocks(interface{}, string, string) error
}
