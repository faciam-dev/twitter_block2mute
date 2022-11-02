package handler

type SessionHandler interface {
	SetContextHandler(ContextHandler)
	Set(string, string)
	Get(string) interface{}
	Delete(string)
	Clear()
	Save() error
}
