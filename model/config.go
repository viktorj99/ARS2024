package model

type Config struct {
	Name       string            `json:"name"`
	Version    int               `json:"version"`
	Parameters map[string]string `json:"params"`
	Labels     map[string]string `json:"labels"`
}

type ConfigRepository interface {
	AddConfig(config Config) error
	GetConfig(name string, version int) (Config, error)
	DeleteConfig(name string, version int) error
}
