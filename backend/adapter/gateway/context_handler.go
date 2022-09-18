package gateway

type ContextHandler interface {
    GetContext() interface{}
    Query(key string) (value string)
    Param(key string) string
    JSON(code int, obj interface{})
}