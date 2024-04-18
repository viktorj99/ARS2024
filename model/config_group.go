package model

type ConfigGroup struct {
	Name           string   `json:"name"`
	Version        int      `json:"version"`
	Configurations []Config `json:"config"`
}

type ConfigGroupRepository interface {
	AddConfigGroup(configGroup ConfigGroup)
	GetConfigGroup(name string, version int) (ConfigGroup, error)
	DeleteConfigGroup(name string, version int) error
	AddConfigToGroup(name string, version int, config Config) error
}
