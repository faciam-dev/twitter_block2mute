package gateway

type SessionHandler interface {
	SetContextHandler(ContextHandler)
	Set(string, string)
	Get(string) interface{}
	Save() error
}
