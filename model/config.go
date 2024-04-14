package model

type Config struct {
	Name       string            `json:"name"`
	Version    int               `json:"version"`
	Parameters map[string]string `json:"params"`
}

type ConfigRepository interface {
	AddConfig(config Config)
	GetConfig(name string, version int) (Config, error)
}
