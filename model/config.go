package model

type Config struct {
	Name       string
	Version    string
	Parameters map[string]string
}

type ConfigRepository interface {
	AddConfig(config Config) error
	DeleteConfig(name string, version string) error
	UpdateConfig(name string, version string, config Config) error
	ListConfigs() ([]Config, error)
	FindConfigByNameAndVersion(name string, version string) (*Config, error)
}
