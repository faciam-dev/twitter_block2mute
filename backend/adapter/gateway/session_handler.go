package gateway

type SessionHandler interface {
	SetContext(interface{})
	Set(string, string)
	Get(string) interface{}
	Save() error
}
