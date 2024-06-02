package service

import (
	"context"
	"projekat/model"

	"go.opentelemetry.io/otel/trace"
)

// Klasa
type ConfigGroupService struct {
	repository model.ConfigGroupRepository
	tracer     trace.Tracer
}

func NewConfigGroupService(repository model.ConfigGroupRepository, tracer trace.Tracer) ConfigGroupService {
	return ConfigGroupService{
		repository: repository,
		tracer:     tracer,
	}
}

func (s ConfigGroupService) AddConfigGroup(ctx context.Context, config model.ConfigGroup) error {
	ctx, span := s.tracer.Start(ctx, "AddConfigGroup")
	defer span.End()
	return s.repository.AddConfigGroup(ctx, config)
}

func (s ConfigGroupService) GetConfigGroup(ctx context.Context, name string, version int) (model.ConfigGroup, error) {
	ctx, span := s.tracer.Start(ctx, "GetConfigGroup")
	defer span.End()
	return s.repository.GetConfigGroup(ctx, name, version)
}

func (s ConfigGroupService) DeleteConfigGroup(ctx context.Context, name string, version int) error {
	ctx, span := s.tracer.Start(ctx, "DeleteConfigGroup")
	defer span.End()
	return s.repository.DeleteConfigGroup(ctx, name, version)
}

func (s ConfigGroupService) AddConfigToGroup(ctx context.Context, name string, version int, config model.Config) error {
	ctx, span := s.tracer.Start(ctx, "AddConfigToGroup")
	defer span.End()
	return s.repository.AddConfigToGroup(ctx, name, version, config)
}

func (s ConfigGroupService) DeleteConfigFromGroup(ctx context.Context, groupName string, groupVersion int, configName string, configVersion int) error {
	ctx, span := s.tracer.Start(ctx, "DeleteConfigFromGroup")
	defer span.End()
	return s.repository.DeleteConfigFromGroup(ctx, groupName, groupVersion, configName, configVersion)
}

func (s ConfigGroupService) GetConfigsFromGroupByLabel(ctx context.Context, groupName string, groupVersion int, labels string) ([]model.Config, error) {
	ctx, span := s.tracer.Start(ctx, "GetConfigsFromGroupByLabel")
	defer span.End()
	return s.repository.GetConfigsFromGroupByLabel(ctx, groupName, groupVersion, labels)
}

func (s ConfigGroupService) DeleteConfigsFromGroupByLabel(ctx context.Context, groupName string, groupVersion int, labels string) error {
	ctx, span := s.tracer.Start(ctx, "DeleteConfigsFromGroupByLabel")
	defer span.End()
	return s.repository.DeleteConfigsFromGroupByLabel(ctx, groupName, groupVersion, labels)
}
