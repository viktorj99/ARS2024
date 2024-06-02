package model

type ConfigGroup struct {
	// Name of the configuration group
	// Required: true
	Name string `json:"name"`

	// Version of the configuration group
	// Required: true
	Version int `json:"version"`

	// Configurations in the group
	// Required: true
	Configurations []Config `json:"configurations"`
}

type ConfigGroupRepository interface {
	AddConfigGroup(configGroup ConfigGroup) error
	GetConfigGroup(name string, version int) (ConfigGroup, error)
	DeleteConfigGroup(name string, version int) error
	AddConfigToGroup(name string, version int, config Config) error
	DeleteConfigFromGroup(groupName string, groupVersion int, configName string, configVersion int) error
	GetConfigsFromGroupByLabel(groupName string, groupVersion int, labels string) ([]Config, error)
	DeleteConfigsFromGroupByLabel(groupName string, groupVersion int, labels string) error
}
