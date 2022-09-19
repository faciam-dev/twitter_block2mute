package handler

// このハンドラは各gatewayで共通の基本的な操作のみを扱う
type DbHandler interface {
	First(interface{}, string) error
	Create(interface{}) error
	Update(interface{}, string) error
	Upsert(interface{}, string, string) error
	Find(interface{}, string, string) error
}

// user独自の処理を追加したハンドラ
type UserDbHandler interface {
	DbHandler
	FindByTwitterID(interface{}, string) error
}