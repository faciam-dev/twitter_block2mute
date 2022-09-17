package gateway

type DbHandler interface {
	First(interface{}, string) error
}
