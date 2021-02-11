package store

type Store interface {
	Set(key string, value string) (string, error)
	Get(key string) (string, error)
}
