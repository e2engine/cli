package environment

import (
	"context"

	"github.com/e2engine/core/model"

	"github.com/e2engine/cli/internal/service"
)

type internalService interface {
	ValidateEnvironment(ctx context.Context, specPath string) (*model.Environment, error)
	ListEnvironments(ctx context.Context, params *service.ListEnvironmentsParams) ([]model.Environment, error)
	CreateEnvironment(ctx context.Context, specPath string) (*model.Environment, error)
	GetEnvironment(ctx context.Context, ref string) (*model.Environment, error)
	DeleteEnvironment(ctx context.Context, ref string) (*model.Environment, error)
}

//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=service.go -destination=mock_service_test.go -package=environment
