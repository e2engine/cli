package test

import (
	"context"

	"github.com/e2engine/core/model"

	"github.com/e2engine/cli/internal/service"
)

type internalService interface {
	ValidateTest(ctx context.Context, specPath string) (*model.Test, error)
	ListTests(ctx context.Context, params *service.ListTestsParams) ([]model.Test, error)
	CreateTest(ctx context.Context, specPath string) (*model.Test, error)
	GetTest(ctx context.Context, ref string) (*model.Test, error)
	DeleteTest(ctx context.Context, ref string) (*model.Test, error)
	RunTest(
		ctx context.Context,
		environmentRef string,
		testRef string,
	) (string, error)
}

//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=service.go -destination=mock_service_test.go -package=test
