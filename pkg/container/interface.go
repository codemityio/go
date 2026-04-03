package container

// Container container service.
type Container interface {
	Get(key string) any
	Set(key string, service any, ss func(service any))
	Keys() []string
}
