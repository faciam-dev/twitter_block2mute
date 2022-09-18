package handler

type DbHandler interface {
	First(interface{}, string) error
}
