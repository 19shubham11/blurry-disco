package store

type Store interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Incr(key string) (int, error)
	Mset(set map[string]string) error
	Mget(keys []string) ([]string, error)
}
