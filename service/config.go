package service

import (
	"context"
	"projekat/model"

	"go.opentelemetry.io/otel/trace"
)

// Klasa
type ConfigService struct {
	repository model.ConfigRepository
	tracer     trace.Tracer
}

func NewConfigService(repository model.ConfigRepository, tracer trace.Tracer) ConfigService {
	return ConfigService{
		repository: repository,
		tracer:     tracer,
	}
}

func (s ConfigService) AddConfig(ctx context.Context, config model.Config) error {
	ctx, span := s.tracer.Start(ctx, "AddConfigService")
	defer span.End()

	return s.repository.AddConfig(ctx, config)
}

func (s ConfigService) GetConfig(ctx context.Context, name string, version int) (model.Config, error) {
	ctx, span := s.tracer.Start(ctx, "GetConfigService")
	defer span.End()

	return s.repository.GetConfig(ctx, name, version)
}

func (s ConfigService) DeleteConfig(ctx context.Context, name string, version int) error {
	ctx, span := s.tracer.Start(ctx, "DeleteConfigService")
	defer span.End()

	return s.repository.DeleteConfig(ctx, name, version)
}
