package config

type Configs interface {
	Open(file string) error
	Get(path ...string) (interface{}, error)
	GetString(path ...string) (string, error)
	GetInt(path ...string) (int64, error)
	GetFloat(path ...string) (float64, error)
	GetSlice(path ...string) ([]string, error)
	GetStrMap(path ...string) (map[string]string, error)
	GetMap(path ...string) (map[string]interface{}, error)
	Watch(path ...string) (Watcher, error)
}
