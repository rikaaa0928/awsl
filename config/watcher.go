package config

type Watcher interface {
	Next() (interface{}, error)
	Stop() error
}
