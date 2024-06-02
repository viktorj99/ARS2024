package model

import "context"

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
	AddConfigGroup(ctx context.Context, configGroup ConfigGroup) error
	GetConfigGroup(ctx context.Context, name string, version int) (ConfigGroup, error)
	DeleteConfigGroup(ctx context.Context, name string, version int) error
	AddConfigToGroup(ctx context.Context, name string, version int, config Config) error
	DeleteConfigFromGroup(ctx context.Context, groupName string, groupVersion int, configName string, configVersion int) error
	GetConfigsFromGroupByLabel(ctx context.Context, groupName string, groupVersion int, labels string) ([]Config, error)
	DeleteConfigsFromGroupByLabel(ctx context.Context, groupName string, groupVersion int, labels string) error
}
