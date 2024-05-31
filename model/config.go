package model

// swagger:model Config
type Config struct {
	// Name of the configuration
	// Required: true
	Name string `json:"name"`

	// Version of the configuration
	// Required: true
	Version int `json:"version"`

	// Parameters are key-value pairs for configuration
	// Required: true
	Parameters map[string]string `json:"params"`

	// Labels are key-value pairs for labeling the configuration
	// Required: true
	Labels map[string]string `json:"labels"`
}

type ConfigRepository interface {
	AddConfig(config Config) error
	GetConfig(name string, version int) (Config, error)
	DeleteConfig(name string, version int) error
}
